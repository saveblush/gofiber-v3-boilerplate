package ftpupload

import (
	"errors"
	"fmt"
	"io"
	"path"
	"strings"
	"time"

	//"github.com/saveblush/ftp"
	"github.com/jlaffaye/ftp"

	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/generic"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/utils/logger"
)

var (
	timeout = 5 * time.Second
)

type client struct {
	*ftp.ServerConn
}

type Auth struct {
	Host     string
	UserName string
	Password string
}

// CcDir change current dir
func (c *client) ensureDir(pathDir string) error {
	if pathDir == "" {
		return errors.New("path is required")
	}

	var joinPath []string
	paths := strings.SplitSeq(pathDir, "/")
	for p := range paths {
		if p == "" {
			continue
		}

		joinPath = append(joinPath, p)
		dir := fmt.Sprintf("/%s", strings.Join(joinPath, "/"))

		//  ไปยัง dir และเป็นการตรวจสอบ dir ว่ามีหรือยัง
		err := c.ChangeDir(dir)
		if err != nil {
			// ถือว่ายังไม่มี dir
			// สร้าง dir
			err = c.MakeDir(dir)
			if err != nil {
				return err
			}

			// กำหนด permission
			// ปิดส่วน กำหนด permission และกลับไปใช้ lib github.com/jlaffaye/ftp
			// err = c.ChangePermission("0777", dir)
			// if err != nil {
			// 	return err
			// }

			// ไปยัง dir
			err = c.ChangeDir(dir)
			if err != nil {
				return err
			}
		}
	}

	// dir
	curDir, err := c.CurrentDir()
	if err != nil {
		return err
	}

	// เช็ค dir ตรงตามที่ต้องการ
	if !generic.Equal(pathDir, curDir) {
		return errors.New("dir not found")
	}

	return nil
}

func Upload(bucketName, prefix, objectName string, object io.Reader, auth *Auth) error {
	// connect
	c, err := ftp.Dial(auth.Host, ftp.DialWithTimeout(timeout))
	if err != nil {
		return err
	}
	defer c.Quit()
	logger.Log.Debug("FTP connect ok")

	// login
	err = c.Login(auth.UserName, auth.Password)
	if err != nil {
		return err
	}
	logger.Log.Debug("FTP login ok")

	if prefix != "" {
		bucketName = path.Join(bucketName, prefix)
	}

	// ไปยัง dir
	err = c.ChangeDir(bucketName)
	if err != nil {
		// กรณีไม่เจอ dir
		client := client{c}
		err = client.ensureDir(bucketName)
		if err != nil {
			return err
		}
	}
	logger.Log.Debug("FTP change dir ok")

	// put
	dst := path.Join(bucketName, objectName)
	logger.Log.Debugf("FTP dst: %s", dst)

	err = c.Stor(dst, object)
	if err != nil {
		return err
	}

	return nil
}
