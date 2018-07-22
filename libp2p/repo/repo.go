package repo

import (
	"crypto/rand"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/dynamicgo/slf4go"

	crypto "github.com/libp2p/go-libp2p-crypto"
)

// Repository .
type Repository interface {
	PrivateKey() (crypto.PrivKey, error)
	PeerStoragePath() (string, error)
}

type repositoryImpl struct {
	sync.RWMutex
	slf4go.Logger
	rootpath string
}

// Open open repository by disk path
func Open(rootpath string) (Repository, error) {

	logger := slf4go.Get("libp2p-repo")

	logger.DebugF("init with root path %s", rootpath)

	fullpath, err := filepath.Abs(rootpath)

	if err != nil {
		return nil, err
	}

	return &repositoryImpl{
		Logger:   logger,
		rootpath: fullpath,
	}, nil
}

func isExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return false
	}

	return true
}

func (impl *repositoryImpl) checkRootPath() error {
	if isExist(impl.rootpath) {
		impl.DebugF("root path exist :%s", impl.rootpath)
		return nil
	}

	impl.DebugF("create root path :%s", impl.rootpath)

	return os.MkdirAll(impl.rootpath, 0777)
}

func (impl *repositoryImpl) PrivateKey() (key crypto.PrivKey, err error) {

	impl.Lock()
	defer impl.Unlock()

	if err := impl.checkRootPath(); err != nil {
		return nil, err
	}

	path := filepath.Join(impl.rootpath, "peer", "peer.key")

	if isExist(path) {
		impl.DebugF("loading private key :%s", path)

		buff, err := ioutil.ReadFile(path)

		if err != nil {
			impl.DebugF("loading private key errr %s", err)
			return nil, err
		}

		impl.DebugF("unmarshal private key ...")

		return crypto.UnmarshalPrivateKey(buff)
	}

	if !isExist(filepath.Dir(path)) {
		if err := os.MkdirAll(filepath.Dir(path), 0777); err != nil {
			return nil, err
		}
	}

	impl.DebugF("generate private key :%s", path)

	key, _, err = crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, rand.Reader)

	if err != nil {
		return key, err
	}

	buff, err := crypto.MarshalPrivateKey(key)

	if err != nil {
		return nil, err
	}

	if err := ioutil.WriteFile(path, buff, 0777); err != nil {
		return nil, err
	}

	return key, nil
}

func (impl *repositoryImpl) PeerStoragePath() (string, error) {

	impl.Lock()
	defer impl.Unlock()

	path := filepath.Join(impl.rootpath, "peer", "storage")

	if err := impl.checkRootPath(); err != nil {
		return "", err
	}

	if err := os.MkdirAll(path, 0777); err != nil {
		return "", err
	}

	return path, nil
}
