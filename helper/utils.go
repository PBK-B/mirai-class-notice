/*
 * @Author: Bin
 * @Date: 2021-09-29
 * @FilePath: /class_notice/helper/utils.go
 */
package helper

import (
	"archive/tar"

	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

// 插件安装包解压缩
func UnzipPlugin(file string) (path string, err error) {
	var path_str = "files/plugins/"
	path_id, err := uuid.NewUUID() // 生成一个路径的 uuid
	path_str = path_str + path_id.String() + "/"
	pathDir := filepath.Join(GetCurrentAbPath(), path_str)
	// err = os.MkdirAll(path_str, os.ModePerm) // 创建一个插件存放文件夹
	if err != nil {
		return
	}

	err = UnTar(pathDir, file)
	if err != nil {
		return
	}

	return path_str, nil
}

// 从 tar 压缩文件中读取文件
func TarReaderFile(name string, path string) (content []byte, err error) {

	// 读取文件
	file, err := os.Open(path)
	if err != nil || file == nil {
		return
	}
	defer file.Close()

	// 通过 tar 创建 tar.Reader
	tr := tar.NewReader(file)

	// 现在已经获得了 tar.Reader 结构了，只需要循环里面的数据读取文件就可以了
	for {
		hdr, err := tr.Next()

		switch {
		case err == io.EOF:
			return nil, nil
		case err != nil:
			return nil, err
		case hdr == nil:
			continue
		}

		// fmt.Printf("读取文件: %v \n", hdr.Name)

		// 根据 header 的 Typeflag 字段，判断文件的类型
		switch hdr.Typeflag {
		case tar.TypeReg: // 读取文件

			// 判断是需要读取的文件
			if name == hdr.Name {
				content, err = ioutil.ReadAll(tr)
				return content, err
			}

		}
	}
}

// 解压 tar 压缩文件
func UnTar(dst string, path string) (err error) {

	// 读取文件
	file, err := os.Open(path)
	if err != nil || file == nil {
		return
	}
	defer file.Close()

	// 通过 file 创建 tar.Reader
	tr := tar.NewReader(file)

	// 现在已经获得了 tar.Reader 结构了，只需要循环里面的数据写入文件就可以了
	for {
		hdr, err := tr.Next()

		switch {
		case err == io.EOF:
			return nil
		case err != nil:
			return err
		case hdr == nil:
			continue
		}

		// 处理下保存路径，将要保存的目录加上 header 中的 Name
		// 这个变量保存的有可能是目录，有可能是文件，所以就叫 FileDir 了……
		if b := Exists(dst); !b {
			// 如果根目录不存在的话创建根目录
			if err := os.MkdirAll(dst, 0775); err != nil {
				return err
			}
		}
		dstFileDir := filepath.Join(dst, hdr.Name)

		// 根据 header 的 Typeflag 字段，判断文件的类型
		switch hdr.Typeflag {
		case tar.TypeDir: // 如果是目录时候，创建目录
			// 判断下目录是否存在，不存在就创建
			if b := Exists(dstFileDir); !b {
				// 使用 MkdirAll 不使用 Mkdir ，就类似 Linux 终端下的 mkdir -p，
				// 可以递归创建每一级目录
				if err := os.MkdirAll(dstFileDir, 0775); err != nil {
					return err
				}
			}
		case tar.TypeReg: // 如果是文件就写入到磁盘
			// 创建一个可以读写的文件，权限就使用 header 中记录的权限
			// 因为操作系统的 FileMode 是 int32 类型的，hdr 中的是 int64，所以转换下
			file, err := os.OpenFile(dstFileDir, os.O_CREATE|os.O_RDWR, os.FileMode(hdr.Mode))
			if err != nil {
				return err
			}
			_, err = io.Copy(file, tr)
			if err != nil {
				return err
			}
			// 将解压结果输出显示
			// fmt.Printf("成功解压： %s , 共处理了 %d 个字符\n", dstFileDir, n)

			// 不要忘记关闭打开的文件，因为它是在 for 循环中，不能使用 defer
			// 如果想使用 defer 就放在一个单独的函数中
			file.Close()
		}
	}
}

// 判断所给路径文件/文件夹是否存在
func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// 判断所给路径是否为文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// 判断所给路径是否为文件
func IsFile(path string) bool {
	return !IsDir(path)
}
