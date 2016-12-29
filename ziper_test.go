package ziper

import(
	"testing"
)

func Test_zip(t *testing.T){
	z := &Ziper{}
	err := z.TarGz("D://迅雷下载","d://tmp.tar.gz")
	if err != nil{
		t.Fatalf("targz file failed,err:%v\n",err)
	}
	err = z.UnTarGz("d://tmp.tar.gz","d://tmp2")
	if err != nil{
		t.Fatalf("untargz file failed,err:%v\n",err)
	}
}

