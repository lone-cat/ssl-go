package main

import (
	"os"
	"ssl/keytype"
	"ssl/managers"
	"ssl/storage"
	"ssl/storage/file"
)

func NewBundleManager[T keytype.Private](
	privateKeyFilename string,
	privateKeyPermissions os.FileMode,
	certificateFilename string,
	certificatePermissions os.FileMode,
	privateKeyAndCertificateFilename string,
	privateKeyAndCertificatePermissions os.FileMode,
	certificateChainFilename string,
	certificateChainPermissions os.FileMode,
	intermediateFilename string,
	intermediatePermissions os.FileMode,
	intermediatePattern string,
	intermediatePatternPermissions os.FileMode,
	allInOneFilename string,
	allInOnePermissions os.FileMode,
) (mgr managers.Bundle[T], err error) {
	var privateKeyStorage,
		certificateStorage,
		privateKeyAndCertificateStorage,
		certificateChainStorage,
		intermediateStorage,
		intermediateMultiStorage,
		allInOneStorage storage.Pem

	if privateKeyFilename != `` {
		privateKeyStorage, err = getPemStorageFromFilenameAndPermissions(privateKeyFilename, privateKeyPermissions)
		if err != nil {
			return
		}
	}

	if certificateFilename != `` {
		certificateStorage, err = getPemStorageFromFilenameAndPermissions(certificateFilename, certificatePermissions)
		if err != nil {
			return
		}
	}

	if privateKeyAndCertificateFilename != `` {
		privateKeyAndCertificateStorage, err = getPemStorageFromFilenameAndPermissions(privateKeyAndCertificateFilename, privateKeyAndCertificatePermissions)
		if err != nil {
			return
		}
	}

	if certificateChainFilename != `` {
		certificateChainStorage, err = getPemStorageFromFilenameAndPermissions(certificateChainFilename, certificateChainPermissions)
		if err != nil {
			return
		}
	}

	if intermediateFilename != `` {
		intermediateStorage, err = getPemStorageFromFilenameAndPermissions(intermediateFilename, intermediatePermissions)
		if err != nil {
			return
		}
	}

	if allInOneFilename != `` {
		allInOneStorage, err = getPemStorageFromFilenameAndPermissions(allInOneFilename, allInOnePermissions)
		if err != nil {
			return
		}
	}

	if intermediatePattern != `` {
		var multiByteStorage storage.ByteMulti
		multiByteStorage, err = file.NewByteMultiFile(intermediatePattern, intermediatePatternPermissions)
		if err != nil {
			return
		}
		intermediateMultiStorage, err = storage.NewPemMultibyte(multiByteStorage)
		if err != nil {
			return
		}
	}

	mgr = managers.NewBundle[T](
		privateKeyStorage,
		certificateStorage,
		privateKeyAndCertificateStorage,
		certificateChainStorage,
		intermediateStorage,
		intermediateMultiStorage,
		allInOneStorage,
	)

	return
}

func getPemStorageFromFilenameAndPermissions(filename string, permissions os.FileMode) (store storage.Pem, err error) {
	byteStorage, err := file.NewByteFile(filename, permissions)
	if err != nil {
		return
	}
	multiByteStorage, err := storage.NewByteSingleFileAdapter(byteStorage)
	if err != nil {
		return
	}
	store, err = storage.NewPemMultibyte(multiByteStorage)

	return
}
