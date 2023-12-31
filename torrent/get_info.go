package torrent

import (
	"context"
	"fmt"
)

type Info struct {
	Name   string `json:"name"`
	Length int64  `json:"length"`
	Files  []File `json:"files"`
}

type File struct {
	DisplayPath    string `json:"displayPath"`
	Length         int64  `json:"length"`
	BytesCompleted int64  `json:"bytesCompleted"`
}

func (h *Handler) GetTorrentInfo(ctx context.Context, infoHashHex string) (Info, error) {
	infoHash, err := infoHashFromHexString(infoHashHex)
	if err != nil {
		return Info{}, fmt.Errorf("parse infohash: %w", err)
	}
	tor, ok := h.torrentClient.Torrent(infoHash)
	if !ok {
		return Info{}, fmt.Errorf("torrent not found")
	}
	<-tor.GotInfo()

	torInfo := Info{
		Name:   tor.Name(),
		Length: tor.Length(),
		Files:  make([]File, 0),
	}

	for _, file := range tor.Files() {
		torInfo.Files = append(torInfo.Files, File{
			DisplayPath:    file.DisplayPath(),
			Length:         file.Length(),
			BytesCompleted: file.BytesCompleted(),
		})
	}
	return torInfo, nil
}
