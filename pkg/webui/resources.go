package webui

import "os"

// Resources Dir Structure:
// {base_path}
//   |__ templates
//   |    |__ webui
//   |    |    |__ ... {webui templates}
//   |    |__ ... {application templates}
//   |__ static
//        |__ webui
//        |    |__ ... {webui static files}
//        |__ vendor
//        |    |__ ... {webui static files}
//        |__ ... {application static files}

var (
	libUrl    = "https://github.com/exonlabs/go-utils"
	libRef    = "heads/master"
	staticUrl = "https://github.com/exonlabs/exonwebui-static"
	staticRef = "tags/v1.0"
)

// https://github.com/exonlabs/go-utils/archive/refs/heads/master.zip

func InitResources(dst_dir string) error {
	os.RemoveAll(dst_dir)

	return nil
}

func InitLibResource() error {
	return nil
}

func InitVendorResource() error {
	return nil
}

// fetch zip archive from web url
func downloadZip(url string, dst_dir string) error {
	return nil
}

func extractZip(src_file string, dst_dir string) error {
	return nil
}
