package scp

import (
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"log"
	"os"
	"scp-copy/config"
)

func BackupDB(sqlConfig *config.SqlConfig, fsConfig *config.FsConfig, bastionConfig *config.BastionConfig) {
	serverSrcPath := fsConfig.DumpDir
	destPath := fsConfig.LocalDir
	fileName := sqlConfig.DBDumpFileName

	cfg := &ssh.ClientConfig{
		User: bastionConfig.User,
		Auth: []ssh.AuthMethod{
			publicKey(fsConfig.PublicKey),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", bastionConfig.Address, bastionConfig.Port), cfg)

	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Successfully connected to Bastion server")
	}

	// open an SFTP session over an existing ssh connection.
	sftpClient, err := sftp.NewClient(conn)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("SFTP client successfully created")
	defer sftpClient.Close()

	session, err := conn.NewSession()

	if err != nil {
		log.Fatal("Error establishing new session: ", err)
	}
	defer session.Close()

	log.Printf("Creating %s database backup", sqlConfig.DBName)

	cmd := getCmd(sqlConfig.Address, sqlConfig.Username, sqlConfig.Password, sqlConfig.DBName, sqlConfig.DBDumpFileName, serverSrcPath)
	if err := session.Run(cmd); err != nil {
		log.Fatal("Error creating DB backup: ", err)
	}

	log.Printf("Copying backup file %s from %s to %s", fileName, serverSrcPath, destPath)

	// Open the source file
	srcFile, err := sftpClient.Open(serverSrcPath + fileName + ".sql")

	if err != nil {
		log.Fatal("Error opening file: ", err)
	}
	defer srcFile.Close()

	// Create the destination file
	dstFile, err := os.Create(destPath + fileName + ".sql")

	if err != nil {
		log.Fatal("Error creating destination backup dir: ", err)
	}
	defer dstFile.Close()

	// Copy the file
	_, err = srcFile.WriteTo(dstFile)

	if err != nil {
		log.Fatal("Error writing file to destination: ", err)
	}

	log.Println("DB backed up successfully")
}

func publicKey(path string) ssh.AuthMethod {
	key, err := os.ReadFile(path)

	if err != nil {
		log.Fatal(err)
	}

	signer, err := ssh.ParsePrivateKey(key)

	if err != nil {
		log.Fatal(err)
	}

	return ssh.PublicKeys(signer)
}

func getCmd(host string, user string, pw string, dbName string, dumpDbName string, dumpPath string) string {
	return fmt.Sprintf(
		"docker exec %s /usr/bin/mysqldump --user=%s --password=%s %s > %s%s.sql",
		host,
		user,
		pw,
		dbName,
		dumpPath,
		dumpDbName,
	)
}
