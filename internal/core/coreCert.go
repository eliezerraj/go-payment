package core

type Cert struct {
	isTLS				bool
	CertPEM 			[]byte
	CertAccountPEM 		[]byte 	 		
	CertPrivKeyPEM	    []byte 
}