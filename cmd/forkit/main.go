package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"github.com/spf13/cobra"
)

const (
	baseName     = "base"
	manifestName = "manifest.json"
	forkMarker   = ".forkit."
)

var (
	outputJSON    bool
	variantNameRE = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_-]*$`)
)

type manifest struct {
	Variants []string `json:"variants"`
}

type paths struct {
	Source   string
	Dir      string
	Stem     string
	Ext      string
	Manifest string
}

type variantInfo struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type listOutput struct {
	Source   string        `json:"source"`
	Base     string        `json:"base"`
	Manifest string        `json:"manifest"`
	Variants []variantInfo `json:"variants"`
}

func main() {
	root := &cobra.Command{
		Use:   "forkit",
		Short: "Embedded implementation prototypes",
		Run: func(cmd *cobra.Command, _ []string) {
			cmd.Help()
		},
	}
	root.PersistentFlags().BoolVar(&outputJSON, "json", false, "Output JSON")
	root.AddCommand(initCmd(), copyCmd(), listCmd(), dropCmd(), cleanCmd(), promoteCmd())
	root.CompletionOptions.DisableDefaultCmd = true

	if err := root.Execute(); err != nil {
		fatalf("%v", err)
	}
}

func initCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "init <file>",
		Short:   "Create sibling prototype files",
		Example: "  forkit init src/components/Button.tsx",
		Args:    cobra.ExactArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			if err := initFile(args[0]); err != nil {
				fatalf("%v", err)
			}
		},
	}
}

func copyCmd() *cobra.Command {
	var from string
	cmd := &cobra.Command{
		Use:     "copy <file> <variant>",
		Aliases: []string{"fork", "dup"},
		Short:   "Copy the source file or another variant",
		Example: "  forkit copy src/components/Button.tsx agent-a\n  forkit copy src/components/Button.tsx compact --from base",
		Args:    cobra.ExactArgs(2),
		Run: func(_ *cobra.Command, args []string) {
			if err := copyVariant(args[0], args[1], from); err != nil {
				fatalf("%v", err)
			}
		},
	}
	cmd.Flags().StringVar(&from, "from", "", "Source variant to copy; defaults to the canonical file")
	return cmd
}

func listCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "list <file>",
		Aliases: []string{"ls"},
		Short:   "List variants for a file",
		Example: "  forkit list src/components/Button.tsx",
		Args:    cobra.ExactArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			p, m, err := load(args[0])
			if err != nil {
				fatalf("%v", err)
			}
			printList(p, m)
		},
	}
}

func dropCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "drop <file> <variant>",
		Short:   "Remove one variant",
		Example: "  forkit drop src/components/Button.tsx agent-a",
		Args:    cobra.ExactArgs(2),
		Run: func(_ *cobra.Command, args []string) {
			if err := dropVariant(args[0], args[1]); err != nil {
				fatalf("%v", err)
			}
		},
	}
}

func cleanCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "clean <file>",
		Short:   "Remove all forkit files for a source file",
		Example: "  forkit clean src/components/Button.tsx",
		Args:    cobra.ExactArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			if err := cleanFile(args[0]); err != nil {
				fatalf("%v", err)
			}
		},
	}
}

func promoteCmd() *cobra.Command {
	var clean bool
	cmd := &cobra.Command{
		Use:     "promote <file> <variant>",
		Short:   "Replace the source file with a variant",
		Example: "  forkit promote src/components/Button.tsx agent-a\n  forkit promote src/components/Button.tsx agent-a --clean",
		Args:    cobra.ExactArgs(2),
		Run: func(_ *cobra.Command, args []string) {
			if err := promoteVariant(args[0], args[1], clean); err != nil {
				fatalf("%v", err)
			}
		},
	}
	cmd.Flags().BoolVar(&clean, "clean", false, "Remove forkit prototype files after promotion")
	return cmd
}

func initFile(raw string) error {
	p := resolve(raw)
	if _, err := os.Stat(p.Source); err != nil {
		return fmt.Errorf("read source: %w", err)
	}
	if _, err := os.Stat(p.Manifest); err == nil {
		return fmt.Errorf("%s is already initialized", p.Source)
	}
	m := manifest{Variants: []string{}}
	if err := copyFile(p.Source, basePath(p)); err != nil {
		return err
	}
	if err := writeManifest(p, m); err != nil {
		return err
	}
	printAction("initialized", p, "")
	return nil
}

func copyVariant(raw, name, from string) error {
	p, m, err := load(raw)
	if err != nil {
		return err
	}
	if err := validateNewVariant(m, name); err != nil {
		return err
	}

	src := p.Source
	if from != "" {
		src, err = variantPath(p, m, from)
		if err != nil {
			return err
		}
	}
	if err := copyFile(src, namedPath(p, name)); err != nil {
		return err
	}

	m.Variants = append(m.Variants, name)
	slices.Sort(m.Variants)
	if err := writeManifest(p, m); err != nil {
		return err
	}
	printAction("copied", p, name)
	return nil
}

func promoteVariant(raw, name string, clean bool) error {
	p, m, err := load(raw)
	if err != nil {
		return err
	}
	src, err := variantPath(p, m, name)
	if err != nil {
		return err
	}
	if err := copyFile(src, p.Source); err != nil {
		return err
	}
	if clean {
		if err := cleanFiles(p, m); err != nil {
			return err
		}
	}
	printAction("promoted", p, name)
	return nil
}

func dropVariant(raw, name string) error {
	p, m, err := load(raw)
	if err != nil {
		return err
	}
	if name == baseName {
		return fmt.Errorf("cannot drop %q; use clean to remove all forkit files", baseName)
	}
	if !slices.Contains(m.Variants, name) {
		return fmt.Errorf("unknown variant %q", name)
	}
	if err := os.Remove(namedPath(p, name)); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("remove variant: %w", err)
	}
	m.Variants = slices.DeleteFunc(m.Variants, func(v string) bool { return v == name })
	if err := writeManifest(p, m); err != nil {
		return err
	}
	printAction("dropped", p, name)
	return nil
}

func cleanFile(raw string) error {
	p, m, err := load(raw)
	if err != nil {
		return err
	}
	if err := cleanFiles(p, m); err != nil {
		return err
	}
	printAction("cleaned", p, "")
	return nil
}

func load(raw string) (paths, manifest, error) {
	p := resolve(raw)
	data, err := os.ReadFile(p.Manifest)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return p, manifest{}, fmt.Errorf("%s is not initialized; run forkit init %s first", p.Source, p.Source)
		}
		return p, manifest{}, fmt.Errorf("read manifest: %w", err)
	}
	var m manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return p, m, fmt.Errorf("parse manifest: %w", err)
	}
	if m.Variants == nil {
		m.Variants = []string{}
	}
	return p, m, nil
}

func writeManifest(p paths, m manifest) error {
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal manifest: %w", err)
	}
	return os.WriteFile(p.Manifest, append(data, '\n'), 0o644)
}

func variantPath(p paths, m manifest, name string) (string, error) {
	if name == baseName {
		return basePath(p), nil
	}
	if slices.Contains(m.Variants, name) {
		return namedPath(p, name), nil
	}
	return "", fmt.Errorf("unknown variant %q", name)
}

func validateNewVariant(m manifest, name string) error {
	if name == baseName {
		return fmt.Errorf("%q is reserved", baseName)
	}
	if err := validateName(name); err != nil {
		return err
	}
	if slices.Contains(m.Variants, name) {
		return fmt.Errorf("variant %q already exists", name)
	}
	return nil
}

func validateName(name string) error {
	if !variantNameRE.MatchString(name) {
		return fmt.Errorf("invalid variant %q; use letters, numbers, dashes, and underscores", name)
	}
	return nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("open %s: %w", src, err)
	}
	defer in.Close()

	info, err := in.Stat()
	if err != nil {
		return fmt.Errorf("stat %s: %w", src, err)
	}

	out, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("create %s: %w", dst, err)
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return fmt.Errorf("copy %s to %s: %w", src, dst, err)
	}
	if err := out.Chmod(info.Mode()); err != nil {
		return fmt.Errorf("chmod %s: %w", dst, err)
	}
	return nil
}

func resolve(raw string) paths {
	source := filepath.Clean(raw)
	dir := filepath.Dir(source)
	base := filepath.Base(source)
	ext := filepath.Ext(base)
	stem := strings.TrimSuffix(base, ext)
	return paths{
		Source:   source,
		Dir:      dir,
		Stem:     stem,
		Ext:      ext,
		Manifest: filepath.Join(dir, stem+forkMarker+manifestName),
	}
}

func basePath(p paths) string {
	return namedPath(p, baseName)
}

func namedPath(p paths, name string) string {
	return filepath.Join(p.Dir, p.Stem+forkMarker+name+p.Ext)
}

func cleanFiles(p paths, m manifest) error {
	files := []string{basePath(p), p.Manifest}
	for _, v := range m.Variants {
		files = append(files, namedPath(p, v))
	}
	for _, file := range files {
		if err := os.Remove(file); err != nil && !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("remove %s: %w", file, err)
		}
	}
	return nil
}

func printList(p paths, m manifest) {
	if outputJSON {
		printJSON(listPayload(p, m))
		return
	}
	fmt.Printf("%s\n", p.Source)
	fmt.Printf("  base\t%s\n", basePath(p))
	for _, v := range m.Variants {
		fmt.Printf("  %s\t%s\n", v, namedPath(p, v))
	}
}

func listPayload(p paths, m manifest) listOutput {
	variants := make([]variantInfo, 0, len(m.Variants))
	for _, v := range m.Variants {
		variants = append(variants, variantInfo{
			Name: v,
			Path: namedPath(p, v),
		})
	}
	return listOutput{
		Source:   p.Source,
		Base:     basePath(p),
		Manifest: p.Manifest,
		Variants: variants,
	}
}

func printAction(action string, p paths, variant string) {
	if outputJSON {
		printJSON(map[string]string{"action": action, "file": p.Source, "variant": variant})
		return
	}
	if variant == "" {
		fmt.Printf("%s %s\n", action, p.Source)
		return
	}
	fmt.Printf("%s %s as %s\n", action, p.Source, variant)
}

func printJSON(v any) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(v); err != nil {
		fatalf("encode JSON: %v", err)
	}
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
