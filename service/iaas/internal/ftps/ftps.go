package ftps

import (
	"crypto/tls"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/jlaffaye/ftp"
	"github.com/sacloud/iaas-api-go"
)

func NewClient(ftpServer *iaas.FTPServer) (*ftp.ServerConn, error) {
	ftpsClient, err := ftp.Dial(
		fmt.Sprintf("%s:%d", ftpServer.HostName, 21),
		ftp.DialWithTimeout(30*time.Minute),
		ftp.DialWithExplicitTLS(&tls.Config{
			ServerName: ftpServer.HostName,
			MinVersion: tls.VersionTLS12,
			MaxVersion: tls.VersionTLS13,
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("cannot create ftp client: %v", err)
	}
	if err := ftpsClient.Login(ftpServer.User, ftpServer.Password); err != nil {
		return nil, fmt.Errorf("cannot login: %v", err)
	}
	return ftpsClient, nil
}

// Download creates a local file at filePath and writes the contents of the first regular file
// (not starting with ".") found in the FTP server's root directory to it.
//
// - Create a local file at filePath
// - Use DownloadWriter to download the first regular file from the FTP server and write to the local file
func Download(ftpClient *ftp.ServerConn, filePath string) error {
	file, err := os.Create(filePath) //nolint:gosec
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	if err := DownloadWriter(ftpClient, file); err != nil {
		return fmt.Errorf("ftps: DownloadWriter: %w", err)
	}
	return nil
}

// DownloadWriter downloads the first regular file (not starting with ".") found in the FTP server's root directory and
// writes it to the given io.Writer.
//
// - List files on the server
// - Find the first regular file (not starting with ".")
// - Download it to the given io.Writer
func DownloadWriter(ftpClient *ftp.ServerConn, writer io.Writer) error {
	entries, err := ftpClient.List("")
	if err != nil {
		return err
	}
	var serverFilePath string
	for _, e := range entries {
		if e.Type == ftp.EntryTypeFile && !strings.HasPrefix(e.Name, ".") {
			serverFilePath = e.Name
			break
		}
	}
	if serverFilePath == "" {
		return fmt.Errorf("FTP retrieve filename failed")
	}
	rc, err := ftpClient.Retr(serverFilePath)
	if err != nil {
		return fmt.Errorf("FTP download file is failed: %w", err)
	}
	defer func() { _ = rc.Close() }()
	_, err = io.Copy(writer, rc)
	if err != nil {
		return fmt.Errorf("FTP download file is failed: %w", err)
	}
	return nil
}
