package file

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCopyFilePreservesMode(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "script.sh")
	require.NoError(t, os.WriteFile(src, []byte("#!/bin/sh\necho hi\n"), 0o755))

	dst := filepath.Join(dir, "out", "script.sh")
	require.NoError(t, Copy(src, dst))

	got, err := os.ReadFile(dst)
	require.NoError(t, err)
	assert.Equal(t, "#!/bin/sh\necho hi\n", string(got))

	info, err := os.Stat(dst)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0o755), info.Mode().Perm())
}

func TestCopyDirectoryTree(t *testing.T) {
	dir := t.TempDir()
	srcDir := filepath.Join(dir, "assets")
	require.NoError(t, os.MkdirAll(filepath.Join(srcDir, "nested"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(srcDir, "a.txt"), []byte("a"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(srcDir, "nested", "b.txt"), []byte("b"), 0o644))

	dstDir := filepath.Join(dir, "copy")
	require.NoError(t, Copy(srcDir, dstDir))

	a, err := os.ReadFile(filepath.Join(dstDir, "a.txt"))
	require.NoError(t, err)
	assert.Equal(t, "a", string(a))

	b, err := os.ReadFile(filepath.Join(dstDir, "nested", "b.txt"))
	require.NoError(t, err)
	assert.Equal(t, "b", string(b))
}
