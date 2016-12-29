package ziper

import(
	"testing"
)

func Test_zip(t *testing.T){
	z := &Targz{}
	err := z.TarGz("D://zip_test","d://tmp.tar.gz")
	if err != nil{
		t.Fatalf("targz file failed,err:%v\n",err)
	}

	//return
	err = z.UnTarGz("d://tmp.tar.gz","d://tmp2")
	if err != nil{
		t.Fatalf("untargz file failed,err:%v\n",err)
	}
}

