package ipfs

import (
	"context"
	"os"

	shell "github.com/ipfs/go-ipfs-api"
)

type UploadResult struct {
	CIDs []string `json:"cids"`
	IVs  [][]byte `json:"ivs"`
	Tags [][]byte `json:"tags"`
}

// UploadChunks encrypts each chunk and uploads it to IPFS, returning the ordered CID list
// (Algorithm 5, lines 5-10; Algorithm 2).
func UploadChunks(ctx context.Context, apiURL string, chunks [][]byte, fileKey []byte) (*UploadResult, error) {
	sh := shell.NewShell(apiURL)
	res := &UploadResult{}
	for i, c := range chunks {
		key, err := DeriveChunkKey(fileKey, i)
		if err != nil { return nil, err }
		enc, err := EncryptChunk(key, c)
		if err != nil { return nil, err }
		cid, err := sh.Add(bytesReader(enc.Data))
		if err != nil { return nil, err }
		res.CIDs = append(res.CIDs, cid)
		res.IVs = append(res.IVs, enc.IV)
		res.Tags = append(res.Tags, enc.Tag)
	}
	return res, nil
}

func DownloadChunk(apiURL, cid string) ([]byte, error) {
	sh := shell.NewShell(apiURL)
	r, err := sh.Cat(cid)
	if err != nil { return nil, err }
	defer r.Close()
	return io.ReadAll(r)
}

// small helpers
func bytesReader(b []byte) *bytesReaderT { return &bytesReaderT{b: b} }

type bytesReaderT struct{ b []byte; i int }
func (r *bytesReaderT) Read(p []byte) (int, error) {
	if r.i >= len(r.b) { return 0, os.ErrClosed }
	n := copy(p, r.b[r.i:])
	r.i += n
	return n, nil
}