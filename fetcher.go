// Copyright 2017 Alexey Slaykovsky
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice,
// this list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice,
// this list of conditions and the following disclaimer in the documentation
// and/or other materials provided with the distribution.
//
// 3. Neither the name of the copyright holder nor the names of its
// contributors may be used to endorse or promote products derived from this
// software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
// ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
// LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
// CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
// SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
// INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
// CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
// ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF
// THE POSSIBILITY OF SUCH DAMAGE.

package fetcher

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
	"fmt"
	pb "gopkg.in/cheggaaa/pb.v1"
)

const (
	FilePrefix string = "gegen-"
)

type HTTPFetcher struct {
	Location   string
	ScratchDir string
}

func (f *HTTPFetcher) MakeURL(fileName string) (string, error) {
	if len(fileName) == 0 {
		return "", fmt.Errorf("Empty file name!")
	}
	return strings.Join([]string{f.Location, fileName}, "/"), nil
}

func (f *HTTPFetcher) WriteFile(url string, file *os.File) (int64, error) {
	res, err := http.Get(url)
	if err != nil {
		return -1, err
	}
	defer res.Body.Close()

	size := res.ContentLength

	bar := pb.New(int(size)).SetUnits(pb.U_BYTES)
	bar.Start()

	writer := io.MultiWriter(file, bar)

	written, err := io.Copy(writer, res.Body)
	if err != nil {
		return -1, err
	}
	if written != size {
		return -1, fmt.Errorf("Written and actual sizes differ!")
	}
	return written, nil

}

func (f *HTTPFetcher) AcquireFile(filePath string) (bool, error) {
	newName := strings.Join(
		[]string{FilePrefix, path.Base(filePath), "."}, "")

	tempFile, err := ioutil.TempFile(f.ScratchDir, newName)

	if err != nil {
		return false, err
	}
	defer tempFile.Close()

	_, err = f.WriteFile(filePath, tempFile)
	if err != nil {
		return false, err
	}

	return true, nil
}
