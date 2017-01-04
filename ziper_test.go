package ziper

import(
	"testing"
)

func Test_zip(t *testing.T){
	z := &Targz{}
	err := z.TarGz("D://迅雷下载","d://迅雷下载.tar.gz")
	if err != nil{
		t.Fatalf("targz file failed,err:%v\n",err)
	}

	//return
	err = z.UnTarGz("d://迅雷下载.tar.gz","d://tmp2")
	if err != nil{
		t.Fatalf("untargz file failed,err:%v\n",err)
	}
}

