package core

type Cert struct {
	IsTLS				bool
	CertPEM 			[]byte
	CertAccountPEM 		[]byte 	 		
	CertPrivKeyPEM	    []byte 
}