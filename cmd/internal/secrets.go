package internal

import (
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/kr/pty"
	"github.com/pkg/errors"
	"github.com/tooploox/oya/pkg/secrets"
	"golang.org/x/crypto/ssh/terminal"
)

func SecretsView(path string, stdout, stderr io.Writer) error {
	output, found, err := secrets.Decrypt(path)
	if err != nil {
		return err
	}
	if !found {
		return errors.Errorf("secret file %q not found", path)
	}
	stdout.Write(output)
	return nil
}

func SecretsEdit(path string, stdout, stderr io.Writer) error {
	viewCmd := secrets.ViewCmd(path)
	if err := terminalRun(viewCmd); err != nil {
		return err
	}
	return nil
}

func SecretsEncrypt(path string, stdout, stderr io.Writer) error {
	if err := secrets.Encrypt(path, path); err != nil {
		return err
	}
	return nil
}

func terminalRun(cmd *exec.Cmd) error {
	// Start the command with a pty.
	ptmx, err := pty.Start(cmd)
	if err != nil {
		return err
	}
	// Make sure to close the pty at the end.
	defer func() { _ = ptmx.Close() }() // Best effort.

	// Handle pty size.
	ch := make(chan os.Signal, 1)
	defer close(ch)
	signal.Notify(ch, syscall.SIGWINCH)
	go func() {
		for range ch {
			if err := pty.InheritSize(os.Stdin, ptmx); err != nil {
				log.Printf("error resizing pty: %s", err)
			}
		}
	}()
	ch <- syscall.SIGWINCH // Initial resize.

	// Set stdin in raw mode.
	oldState, err := terminal.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}
	defer func() { _ = terminal.Restore(int(os.Stdin.Fd()), oldState) }() // Best effort.

	// Copy stdin to the pty and the pty to stdout.
	go func() {
		_, err = io.Copy(ptmx, os.Stdin)
		if err != nil {
			log.Printf("error %s", err)
		}
	}()

	_, err = io.Copy(os.Stdout, ptmx)
	if err != nil {
		return err
	}

	return nil
}
