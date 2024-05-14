package ffmpeg

import (
	"log/slog"
	"os"

	"github.com/google/uuid"

	sh "github.com/Watch2Gather/server/internal/pkg/shell"
)

var movieDir = os.Getenv("MOVIE_PATH_PREFIX")

func SplitIntoChunks(inputFile, outputFile, segmentTime string, roomID uuid.UUID) {
	inputFile = movieDir + inputFile
	outputFile = movieDir + "rooms/" + roomID.String() + "/" + outputFile

	sh.MkDir(movieDir+"rooms/", roomID.String())
	args := `-i ` + inputFile + ` -c copy -map 0 -segment_time ` + segmentTime + ` -f segment ` + outputFile + " -loglevel error -f mpegts -codec:v mpeg1video -codec:a mp2 -b 0"
	_, stderr, err := sh.Exec("ffmpeg", args)
	slog.Debug("ffmpeg", "err", stderr, "err_go", err)
	slog.Debug("ff", "cmd", "ffmpeg "+args)
}
