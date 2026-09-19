package automation

import (
	"os"
	"path/filepath"
	"testing"

	"filippo.io/age"
)

func TestBackupFileRegexMatchesPlaintextAndEncrypted(t *testing.T) {
	for _, name := range []string{"openschool_20260101_020000.dump", "openschool_20260101_020000.dump.age"} {
		if !backupFileRe.MatchString(name) {
			t.Fatalf("backupFileRe did not match %q", name)
		}
	}
	if backupFileRe.MatchString("not-a-backup.txt") {
		t.Fatal("backupFileRe matched an unrelated filename")
	}
}

func TestEncryptBackupFileRoundTrips(t *testing.T) {
	identity, err := age.GenerateX25519Identity()
	if err != nil {
		t.Fatal(err)
	}

	dir := t.TempDir()
	plainPath := filepath.Join(dir, "openschool_20260101_020000.dump")
	want := []byte("pretend this is a pg_dump")
	if err := os.WriteFile(plainPath, want, 0o600); err != nil {
		t.Fatal(err)
	}

	encPath, err := encryptBackupFile(plainPath, identity.Recipient().String())
	if err != nil {
		t.Fatal(err)
	}
	if encPath != plainPath+".age" {
		t.Fatalf("encPath = %q", encPath)
	}

	encrypted, err := os.ReadFile(encPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(encrypted) == string(want) {
		t.Fatal("encryptBackupFile wrote the plaintext unchanged")
	}

	f, err := os.Open(encPath)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	r, err := age.Decrypt(f, identity)
	if err != nil {
		t.Fatal(err)
	}
	got := make([]byte, len(want))
	if _, err := r.Read(got); err != nil {
		t.Fatal(err)
	}
	if string(got) != string(want) {
		t.Fatalf("decrypted = %q, want %q", got, want)
	}
}

func TestEncryptBackupFileRejectsInvalidRecipient(t *testing.T) {
	dir := t.TempDir()
	plainPath := filepath.Join(dir, "openschool_20260101_020000.dump")
	if err := os.WriteFile(plainPath, []byte("data"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := encryptBackupFile(plainPath, "not-a-recipient"); err == nil {
		t.Fatal("encryptBackupFile accepted an invalid recipient")
	}
}

func TestCopyBackupFileRejectsSameDirectoryAsSource(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "openschool_20260101_020000.dump")
	want := []byte("dump contents")
	if err := os.WriteFile(src, want, 0o600); err != nil {
		t.Fatal(err)
	}

	if err := copyBackupFile(src, dir); err == nil {
		t.Fatal("copyBackupFile accepted a destination equal to the source directory")
	}

	got, err := os.ReadFile(src)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(want) {
		t.Fatal("copyBackupFile truncated the source file")
	}
}

func TestCopyBackupFileCopiesToDestination(t *testing.T) {
	srcDir, destDir := t.TempDir(), filepath.Join(t.TempDir(), "offsite")
	src := filepath.Join(srcDir, "openschool_20260101_020000.dump")
	want := []byte("dump contents")
	if err := os.WriteFile(src, want, 0o600); err != nil {
		t.Fatal(err)
	}

	if err := copyBackupFile(src, destDir); err != nil {
		t.Fatal(err)
	}

	got, err := os.ReadFile(filepath.Join(destDir, "openschool_20260101_020000.dump"))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(want) {
		t.Fatalf("copied contents = %q, want %q", got, want)
	}
}
