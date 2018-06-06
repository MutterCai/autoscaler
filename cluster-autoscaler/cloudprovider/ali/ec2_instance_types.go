package ali

type instanceType struct {
	InstanceType string
	VCPU         int64
	MemoryMb     int64
	GPU          int64
}

// InstanceTypes is a map of ec2 resources
var InstanceTypes = map[string]*instanceType{
	"ecs.g5.large":
{
	InstanceType: "ecs.g5.large",
	VCPU: 2,
	MemoryMb: 8192,
	GPU: 0,
},
	"ecs.g5.xlarge":
{
	InstanceType: "ecs.g5.xlarge",
	VCPU: 4,
	MemoryMb: 16384,
	GPU: 0,
},
	"ecs.g5.2xlarge":
{
	InstanceType: "ecs.g5.2xlarge",
	VCPU: 8,
	MemoryMb: 32768,
	GPU: 0,
},
	"ecs.g5.4xlarge":
{
	InstanceType: "ecs.g5.4xlarge",
	VCPU: 16,
	MemoryMb: 65536,
	GPU: 0,
},
	"ecs.g5.6xlarge":
{
	InstanceType: "ecs.g5.6xlarge",
	VCPU: 24,
	MemoryMb: 98304,
	GPU: 0,
},
	"ecs.g5.8xlarge":
{
	InstanceType: "ecs.g5.8xlarge",
	VCPU: 32,
	MemoryMb: 131072,
	GPU: 0,
},
	"ecs.g5.16xlarge":
{
	InstanceType: "ecs.g5.16xlarge",
	VCPU: 64,
	MemoryMb: 262144,
	GPU: 0,
},
	"ecs.sn2ne.large":
{
	InstanceType: "ecs.sn2ne.large",
	VCPU: 2,
	MemoryMb: 8192,
	GPU: 0,
},
	"ecs.sn2ne.xlarge":
{
	InstanceType: "ecs.sn2ne.xlarge",
	VCPU: 4,
	MemoryMb: 16384,
	GPU: 0,
},
	"ecs.sn2ne.2xlarge":
{
	InstanceType: "ecs.sn2ne.2xlarge",
	VCPU: 8,
	MemoryMb: 32768,
	GPU: 0,
},
	"ecs.sn2ne.4xlarge":
{
	InstanceType: "ecs.sn2ne.4xlarge",
	VCPU: 16,
	MemoryMb: 65536,
	GPU: 0,
},
	"ecs.sn2ne.8xlarge":
{
	InstanceType: "ecs.sn2ne.8xlarge",
	VCPU: 32,
	MemoryMb: 131072,
	GPU: 0,
},
	"ecs.sn2ne.14xlarge":
{
	InstanceType: "ecs.sn2ne.14xlarge",
	VCPU: 56,
	MemoryMb: 229376,
	GPU: 0,
},
	"ecs.c5.large":
{
	InstanceType: "ecs.c5.large",
	VCPU: 2,
	MemoryMb: 4096,
	GPU: 0,
},
	"ecs.c5.xlarge":
{
	InstanceType: "ecs.c5.xlarge",
	VCPU: 4,
	MemoryMb: 8192,
	GPU: 0,
},
	"ecs.c5.2xlarge":
{
	InstanceType: "ecs.c5.2xlarge",
	VCPU: 8,
	MemoryMb: 16384,
	GPU: 0,
},
	"ecs.c5.4xlarge":
{
	InstanceType: "ecs.c5.4xlarge",
	VCPU: 16,
	MemoryMb: 32768,
	GPU: 0,
},
	"ecs.c5.6xlarge":
{
	InstanceType: "ecs.c5.6xlarge",
	VCPU: 24,
	MemoryMb: 49152,
	GPU: 0,
},
	"ecs.c5.8xlarge":
{
	InstanceType: "ecs.c5.8xlarge",
	VCPU: 32,
	MemoryMb: 65536,
	GPU: 0,
},
	"ecs.c5.16xlarge":
{
	InstanceType: "ecs.c5.16xlarge",
	VCPU: 64,
	MemoryMb: 131072,
	GPU: 0,
},
	"ecs.sn1ne.large":
{
	InstanceType: "ecs.sn1ne.large",
	VCPU: 2,
	MemoryMb: 4096,
	GPU: 0,
},
	"ecs.sn1ne.xlarge":
{
	InstanceType: "ecs.sn1ne.xlarge",
	VCPU: 4,
	MemoryMb: 8192,
	GPU: 0,
},
	"ecs.sn1ne.2xlarge":
{
	InstanceType: "ecs.sn1ne.2xlarge",
	VCPU: 8,
	MemoryMb: 16384,
	GPU: 0,
},
	"ecs.sn1ne.4xlarge":
{
	InstanceType: "ecs.sn1ne.4xlarge",
	VCPU: 16,
	MemoryMb: 32768,
	GPU: 0,
},
	"ecs.sn1ne.8xlarge":
{
	InstanceType: "ecs.sn1ne.8xlarge",
	VCPU: 32,
	MemoryMb: 65536,
	GPU: 0,
},
	"ecs.r5.large":
{
	InstanceType: "ecs.r5.large",
	VCPU: 2,
	MemoryMb: 16384,
	GPU: 0,
},
	"ecs.r5.xlarge":
{
	InstanceType: "ecs.r5.xlarge",
	VCPU: 4,
	MemoryMb: 32768,
	GPU: 0,
},
	"ecs.r5.2xlarge":
{
	InstanceType: "ecs.r5.2xlarge",
	VCPU: 8,
	MemoryMb: 65536,
	GPU: 0,
},
	"ecs.r5.4xlarge":
{
	InstanceType: "ecs.r5.4xlarge",
	VCPU: 16,
	MemoryMb: 131072,
	GPU: 0,
},
	"ecs.r5.6xlarge":
{
	InstanceType: "ecs.r5.6xlarge",
	VCPU: 24,
	MemoryMb: 196608,
	GPU: 0,
},
	"ecs.r5.8xlarge":
{
	InstanceType: "ecs.r5.8xlarge",
	VCPU: 32,
	MemoryMb: 262144,
	GPU: 0,
},
	"ecs.r5.16xlarge":
{
	InstanceType: "ecs.r5.16xlarge",
	VCPU: 64,
	MemoryMb: 524288,
	GPU: 0,
},
	"ecs.re4.20xlarge":
{
	InstanceType: "ecs.re4.20xlarge",
	VCPU: 80,
	MemoryMb: 983040,
	GPU: 0,
},
	"ecs.re4.40xlarge":
{
	InstanceType: "ecs.re4.40xlarge",
	VCPU: 160,
	MemoryMb: 1966080,
	GPU: 0,
},
	"ecs.se1ne.large":
{
	InstanceType: "ecs.se1ne.large",
	VCPU: 2,
	MemoryMb: 16384,
	GPU: 0,
},
	"ecs.se1ne.xlarge":
{
	InstanceType: "ecs.se1ne.xlarge",
	VCPU: 4,
	MemoryMb: 32768,
	GPU: 0,
},
	"ecs.se1ne.2xlarge":
{
	InstanceType: "ecs.se1ne.2xlarge",
	VCPU: 8,
	MemoryMb: 65536,
	GPU: 0,
},
	"ecs.se1ne.4xlarge":
{
	InstanceType: "ecs.se1ne.4xlarge",
	VCPU: 16,
	MemoryMb: 131072,
	GPU: 0,
},
	"ecs.se1ne.8xlarge":
{
	InstanceType: "ecs.se1ne.8xlarge",
	VCPU: 32,
	MemoryMb: 262144,
	GPU: 0,
},
	"ecs.se1ne.14xlarge":
{
	InstanceType: "ecs.se1ne.14xlarge",
	VCPU: 56,
	MemoryMb: 491520,
	GPU: 0,
},
	"ecs.se1.large":
{
	InstanceType: "ecs.se1.large",
	VCPU: 2,
	MemoryMb: 16384,
	GPU: 0,
},
	"ecs.se1.xlarge":
{
	InstanceType: "ecs.se1.xlarge",
	VCPU: 4,
	MemoryMb: 32768,
	GPU: 0,
},
	"ecs.se1.2xlarge":
{
	InstanceType: "ecs.se1.2xlarge",
	VCPU: 8,
	MemoryMb: 65536,
	GPU: 0,
},
	"ecs.se1.4xlarge":
{
	InstanceType: "ecs.se1.4xlarge",
	VCPU: 16,
	MemoryMb: 131072,
	GPU: 0,
},
	"ecs.se1.8xlarge":
{
	InstanceType: "ecs.se1.8xlarge",
	VCPU: 32,
	MemoryMb: 262144,
	GPU: 0,
},
	"ecs.se1.14xlarge":
{
	InstanceType: "ecs.se1.14xlarge",
	VCPU: 56,
	MemoryMb: 491520,
	GPU: 0,
},
	"ecs.d1ne.2xlarge":
{
	InstanceType: "ecs.d1ne.2xlarge",
	VCPU: 8,
	MemoryMb: 32768,
	GPU: 0,
},
	"ecs.d1ne.4xlarge":
{
	InstanceType: "ecs.d1ne.4xlarge",
	VCPU: 16,
	MemoryMb: 65536,
	GPU: 0,
},
	"ecs.d1ne.6xlarge":
{
	InstanceType: "ecs.d1ne.6xlarge",
	VCPU: 24,
	MemoryMb: 98304,
	GPU: 0,
},
	"ecs.d1ne.8xlarge":
{
	InstanceType: "ecs.d1ne.8xlarge",
	VCPU: 32,
	MemoryMb: 131072,
	GPU: 0,
},
	"ecs.d1ne.14xlarge":
{
	InstanceType: "ecs.d1ne.14xlarge",
	VCPU: 56,
	MemoryMb: 229376,
	GPU: 0,
},
	"ecs.d1.2xlarge":
{
	InstanceType: "ecs.d1.2xlarge",
	VCPU: 8,
	MemoryMb: 32768,
	GPU: 0,
},
	"ecs.d1.4xlarge":
{
	InstanceType: "ecs.d1.4xlarge",
	VCPU: 16,
	MemoryMb: 65536,
	GPU: 0,
},
	"ecs.d1.6xlarge":
{
	InstanceType: "ecs.d1.6xlarge",
	VCPU: 24,
	MemoryMb: 98304,
	GPU: 0,
},
	"ecs.d1-c8d3.8xlarge":
{
	InstanceType: "ecs.d1-c8d3.8xlarge",
	VCPU: 32,
	MemoryMb: 131072,
	GPU: 0,
},
	"ecs.d1.8xlarge":
{
	InstanceType: "ecs.d1.8xlarge",
	VCPU: 32,
	MemoryMb: 131072,
	GPU: 0,
},
	"ecs.d1-c14d3.14xlarge":
{
	InstanceType: "ecs.d1-c14d3.14xlarge",
	VCPU: 56,
	MemoryMb: 163840,
	GPU: 0,
},
	"ecs.d1.14xlarge":
{
	InstanceType: "ecs.d1.14xlarge",
	VCPU: 56,
	MemoryMb: 229376,
	GPU: 0,
},
	"ecs.i2.xlarge":
{
	InstanceType: "ecs.i2.xlarge",
	VCPU: 4,
	MemoryMb: 32768,
	GPU: 0,
},
	"ecs.i2.2xlarge":
{
	InstanceType: "ecs.i2.2xlarge",
	VCPU: 8,
	MemoryMb: 65536,
	GPU: 0,
},
	"ecs.i2.4xlarge":
{
	InstanceType: "ecs.i2.4xlarge",
	VCPU: 16,
	MemoryMb: 131072,
	GPU: 0,
},
	"ecs.i2.8xlarge":
{
	InstanceType: "ecs.i2.8xlarge",
	VCPU: 32,
	MemoryMb: 262144,
	GPU: 0,
},
	"ecs.i2.16xlarge":
{
	InstanceType: "ecs.i2.16xlarge",
	VCPU: 64,
	MemoryMb: 524288,
	GPU: 0,
},
	"ecs.i1.xlarge":
{
	InstanceType: "ecs.i1.xlarge",
	VCPU: 4,
	MemoryMb: 16384,
	GPU: 0,
},
	"ecs.i1.2xlarge":
{
	InstanceType: "ecs.i1.2xlarge",
	VCPU: 8,
	MemoryMb: 32768,
	GPU: 0,
},
	"ecs.i1.4xlarge":
{
	InstanceType: "ecs.i1.4xlarge",
	VCPU: 16,
	MemoryMb: 65536,
	GPU: 0,
},
	"ecs.i1-c5d1.4xlarge":
{
	InstanceType: "ecs.i1-c5d1.4xlarge",
	VCPU: 16,
	MemoryMb: 65536,
	GPU: 0,
},
	"ecs.i1.8xlarge":
{
	InstanceType: "ecs.i1.8xlarge",
	VCPU: 32,
	MemoryMb: 131072,
	GPU: 0,
},
	"ecs.i1-c10d1.8xlarge":
{
	InstanceType: "ecs.i1-c10d1.8xlarge",
	VCPU: 32,
	MemoryMb: 131072,
	GPU: 0,
},
	"ecs.i1.14xlarge":
{
	InstanceType: "ecs.i1.14xlarge",
	VCPU: 56,
	MemoryMb: 229376,
	GPU: 0,
},
	"ecs.hfc5.large":
{
	InstanceType: "ecs.hfc5.large",
	VCPU: 2,
	MemoryMb: 4096,
	GPU: 0,
},
	"ecs.hfc5.xlarge":
{
	InstanceType: "ecs.hfc5.xlarge",
	VCPU: 4,
	MemoryMb: 8192,
	GPU: 0,
},
	"ecs.hfc5.2xlarge":
{
	InstanceType: "ecs.hfc5.2xlarge",
	VCPU: 8,
	MemoryMb: 16384,
	GPU: 0,
},
	"ecs.hfc5.4xlarge":
{
	InstanceType: "ecs.hfc5.4xlarge",
	VCPU: 16,
	MemoryMb: 32768,
	GPU: 0,
},
	"ecs.hfc5.6xlarge":
{
	InstanceType: "ecs.hfc5.6xlarge",
	VCPU: 24,
	MemoryMb: 49152,
	GPU: 0,
},
	"ecs.hfc5.8xlarge":
{
	InstanceType: "ecs.hfc5.8xlarge",
	VCPU: 32,
	MemoryMb: 65536,
	GPU: 0,
},
	"ecs.hfg5.large":
{
	InstanceType: "ecs.hfg5.large",
	VCPU: 2,
	MemoryMb: 8192,
	GPU: 0,
},
	"ecs.hfg5.xlarge":
{
	InstanceType: "ecs.hfg5.xlarge",
	VCPU: 4,
	MemoryMb: 16384,
	GPU: 0,
},
	"ecs.hfg5.2xlarge":
{
	InstanceType: "ecs.hfg5.2xlarge",
	VCPU: 8,
	MemoryMb: 32768,
	GPU: 0,
},
	"ecs.hfg5.4xlarge":
{
	InstanceType: "ecs.hfg5.4xlarge",
	VCPU: 16,
	MemoryMb: 65536,
	GPU: 0,
},
	"ecs.hfg5.6xlarge":
{
	InstanceType: "ecs.hfg5.6xlarge",
	VCPU: 24,
	MemoryMb: 98304,
	GPU: 0,
},
	"ecs.hfg5.8xlarge":
{
	InstanceType: "ecs.hfg5.8xlarge",
	VCPU: 32,
	MemoryMb: 131072,
	GPU: 0,
},
	"ecs.hfg5.14xlarge":
{
	InstanceType: "ecs.hfg5.14xlarge",
	VCPU: 56,
	MemoryMb: 163840,
	GPU: 0,
},
	"ecs.c4.xlarge":
{
	InstanceType: "ecs.c4.xlarge",
	VCPU: 4,
	MemoryMb: 8192,
	GPU: 0,
},
	"ecs.c4.2xlarge":
{
	InstanceType: "ecs.c4.2xlarge",
	VCPU: 8,
	MemoryMb: 16384,
	GPU: 0,
},
	"ecs.c4.4xlarge":
{
	InstanceType: "ecs.c4.4xlarge",
	VCPU: 16,
	MemoryMb: 32768,
	GPU: 0,
},
	"ecs.cm4.xlarge":
{
	InstanceType: "ecs.cm4.xlarge",
	VCPU: 4,
	MemoryMb: 16384,
	GPU: 0,
},
	"ecs.cm4.2xlarge":
{
	InstanceType: "ecs.cm4.2xlarge",
	VCPU: 8,
	MemoryMb: 32768,
	GPU: 0,
},
	"ecs.cm4.4xlarge":
{
	InstanceType: "ecs.cm4.4xlarge",
	VCPU: 16,
	MemoryMb: 65536,
	GPU: 0,
},
	"ecs.cm4.6xlarge":
{
	InstanceType: "ecs.cm4.6xlarge",
	VCPU: 24,
	MemoryMb: 98304,
	GPU: 0,
},
	"ecs.ce4.xlarge":
{
	InstanceType: "ecs.ce4.xlarge",
	VCPU: 4,
	MemoryMb: 32768,
	GPU: 0,
},
	"ecs.gn5-c4g1.xlarge":
{
	InstanceType: "ecs.gn5-c4g1.xlarge",
	VCPU: 4,
	MemoryMb: 30720,
	GPU: 1,
    },
	"ecs.gn5-c8g1.2xlarge":
{
	InstanceType: "ecs.gn5-c8g1.2xlarge",
	VCPU: 8,
	MemoryMb: 61440,
	GPU: 1,
    },
	"ecs.gn5-c4g1.2xlarge":
{
	InstanceType: "ecs.gn5-c4g1.2xlarge",
	VCPU: 8,
	MemoryMb: 61440,
	GPU: 2,
    },
	"ecs.gn5-c8g1.4xlarge":
{
	InstanceType: "ecs.gn5-c8g1.4xlarge",
	VCPU: 16,
	MemoryMb: 122880,
	GPU: 2,
    },
	"ecs.gn5-c28g1.7xlarge":
{
	InstanceType: "ecs.gn5-c28g1.7xlarge",
	VCPU: 28,
	MemoryMb: 114688,
	GPU: 1,
    },
	"ecs.gn5-c8g1.8xlarge":
{
	InstanceType: "ecs.gn5-c8g1.8xlarge",
	VCPU: 32,
	MemoryMb: 245760,
	GPU: 4,
    },
	"ecs.gn5-c28g1.14xlarge":
{
	InstanceType: "ecs.gn5-c28g1.14xlarge",
	VCPU: 56,
	MemoryMb: 229376,
	GPU: 2,
    },
	"ecs.gn5-c8g1.14xlarge":
{
	InstanceType: "ecs.gn5-c8g1.14xlarge",
	VCPU: 54,
	MemoryMb: 491520,
	GPU: 8,
    },
	"ecs.gn5i-c2g1.large":
{
	InstanceType: "ecs.gn5i-c2g1.large",
	VCPU: 2,
	MemoryMb: 8192,
	GPU: 1,
    },
	"ecs.gn5i-c4g1.xlarge":
{
	InstanceType: "ecs.gn5i-c4g1.xlarge",
	VCPU: 4,
	MemoryMb: 16384,
	GPU: 1,
    },
	"ecs.gn5i-c8g1.2xlarge":
{
	InstanceType: "ecs.gn5i-c8g1.2xlarge",
	VCPU: 8,
	MemoryMb: 32768,
	GPU: 1,
    },
	"ecs.gn5i-c16g1.4xlarge":
{
	InstanceType: "ecs.gn5i-c16g1.4xlarge",
	VCPU: 16,
	MemoryMb: 65536,
	GPU: 1,
    },
	"ecs.gn5i-c28g1.14xlarge":
{
	InstanceType: "ecs.gn5i-c28g1.14xlarge",
	VCPU: 56,
	MemoryMb: 229376,
	GPU: 2,
    },
	"ecs.gn4-c4g1.xlarge":
{
	InstanceType: "ecs.gn4-c4g1.xlarge",
	VCPU: 4,
	MemoryMb: 30720,
	GPU: 1,
    },
	"ecs.gn4-c8g1.2xlarge":
{
	InstanceType: "ecs.gn4-c8g1.2xlarge",
	VCPU: 8,
	MemoryMb: 61440,
	GPU: 1,
    },
	"ecs.gn4.8xlarge":
{
	InstanceType: "ecs.gn4.8xlarge",
	VCPU: 32,
	MemoryMb: 49152,
	GPU: 1,
    },
	"ecs.gn4-c4g1.2xlarge":
{
	InstanceType: "ecs.gn4-c4g1.2xlarge",
	VCPU: 8,
	MemoryMb: 61440,
	GPU: 2,
    },
	"ecs.gn4-c8g1.4xlarge":
{
	InstanceType: "ecs.gn4-c8g1.4xlarge",
	VCPU: 16,
	MemoryMb: 61440,
	GPU: 2,
    },
	"ecs.gn4.14xlarge":
{
	InstanceType: "ecs.gn4.14xlarge",
	VCPU: 56,
	MemoryMb: 98304,
	GPU: 2,
    },
	"ecs.ga1.xlarge":
{
	InstanceType: "ecs.ga1.xlarge",
	VCPU: 4,
	MemoryMb: 10240,
	GPU: 0,
},
	"ecs.ga1.2xlarge":
{
	InstanceType: "ecs.ga1.2xlarge",
	VCPU: 8,
	MemoryMb: 20480,
	GPU: 0,
},
	"ecs.ga1.4xlarge":
{
	InstanceType: "ecs.ga1.4xlarge",
	VCPU: 16,
	MemoryMb: 40960,
	GPU: 0,
},
	"ecs.ga1.8xlarge":
{
	InstanceType: "ecs.ga1.8xlarge",
	VCPU: 32,
	MemoryMb: 81920,
	GPU: 0,
},
	"ecs.ga1.14xlarge":
{
	InstanceType: "ecs.ga1.14xlarge",
	VCPU: 56,
	MemoryMb: 163840,
	GPU: 0,
},
	"ecs.f1-c8f1.2xlarge":
{
	InstanceType: "ecs.f1-c8f1.2xlarge",
	VCPU: 8,
	MemoryMb: 61440,
	GPU: 0,
},
	"ecs.f1-c28f1.7xlarge":
{
	InstanceType: "ecs.f1-c28f1.7xlarge",
	VCPU: 28,
	MemoryMb: 114688,
	GPU: 0,
},
	"ecs.f2-c8f1.2xlarge":
{
	InstanceType: "ecs.f2-c8f1.2xlarge",
	VCPU: 8,
	MemoryMb: 61440,
	GPU: 0,
},
	"ecs.f2-c8f1.4xlarge":
{
	InstanceType: "ecs.f2-c8f1.4xlarge",
	VCPU: 16,
	MemoryMb: 122880,
	GPU: 0,
},
	"ecs.f2-c28f1.7xlarge":
{
	InstanceType: "ecs.f2-c28f1.7xlarge",
	VCPU: 28,
	MemoryMb: 114688,
	GPU: 0,
},
	"ecs.f2-c28f1.14xlarge":
{
	InstanceType: "ecs.f2-c28f1.14xlarge",
	VCPU: 56,
	MemoryMb: 229376,
	GPU: 0,
},
	"ecs.ebmg5.24xlarge":
{
	InstanceType: "ecs.ebmg5.24xlarge",
	VCPU: 96,
	MemoryMb: 393216,
	GPU: 0,
},
	"ecs.ebmg4.8xlarge":
{
	InstanceType: "ecs.ebmg4.8xlarge",
	VCPU: 32,
	MemoryMb: 131072,
	GPU: 0,
},
	"ecs.ebmhfg5.2xlarge":
{
	InstanceType: "ecs.ebmhfg5.2xlarge",
	VCPU: 8,
	MemoryMb: 32768,
	GPU: 0,
},
	"ecs.ebmhfg4.4xlarge":
{
	InstanceType: "ecs.ebmhfg4.4xlarge",
	VCPU: 16,
	MemoryMb: 65536,
	GPU: 0,
},
	"ecs.ebmc4.8xlarge":
{
	InstanceType: "ecs.ebmc4.8xlarge",
	VCPU: 32,
	MemoryMb: 65536,
	GPU: 0,
},
	"ecs.sccg5.24xlarge":
{
	InstanceType: "ecs.sccg5.24xlarge",
	VCPU: 96,
	MemoryMb: 393216,
	GPU: 0,
},
	"ecs.scch5.16xlarge":
{
	InstanceType: "ecs.scch5.16xlarge",
	VCPU: 64,
	MemoryMb: 196608,
	GPU: 0,
},
	"ecs.t5-lc2m1.nano":
{
	InstanceType: "ecs.t5-lc2m1.nano",
	VCPU: 1,
	MemoryMb: 0,
	GPU: 0,
},
	"ecs.t5-lc1m1.small":
{
	InstanceType: "ecs.t5-lc1m1.small",
	VCPU: 1,
	MemoryMb: 1024,
	GPU: 0,
},
	"ecs.t5-lc1m2.small":
{
	InstanceType: "ecs.t5-lc1m2.small",
	VCPU: 1,
	MemoryMb: 2048,
	GPU: 0,
},
	"ecs.t5-lc1m2.large":
{
	InstanceType: "ecs.t5-lc1m2.large",
	VCPU: 2,
	MemoryMb: 4096,
	GPU: 0,
},
	"ecs.t5-lc1m4.large":
{
	InstanceType: "ecs.t5-lc1m4.large",
	VCPU: 2,
	MemoryMb: 8192,
	GPU: 0,
},
	"ecs.t5-c1m1.large":
{
	InstanceType: "ecs.t5-c1m1.large",
	VCPU: 2,
	MemoryMb: 2048,
	GPU: 0,
},
	"ecs.t5-c1m2.large":
{
	InstanceType: "ecs.t5-c1m2.large",
	VCPU: 2,
	MemoryMb: 4096,
	GPU: 0,
},
	"ecs.t5-c1m4.large":
{
	InstanceType: "ecs.t5-c1m4.large",
	VCPU: 2,
	MemoryMb: 8192,
	GPU: 0,
},
	"ecs.t5-c1m1.xlarge":
{
	InstanceType: "ecs.t5-c1m1.xlarge",
	VCPU: 4,
	MemoryMb: 4096,
	GPU: 0,
},
	"ecs.t5-c1m2.xlarge":
{
	InstanceType: "ecs.t5-c1m2.xlarge",
	VCPU: 4,
	MemoryMb: 8192,
	GPU: 0,
},
	"ecs.t5-c1m4.xlarge":
{
	InstanceType: "ecs.t5-c1m4.xlarge",
	VCPU: 4,
	MemoryMb: 16384,
	GPU: 0,
},
	"ecs.t5-c1m1.2xlarge":
{
	InstanceType: "ecs.t5-c1m1.2xlarge",
	VCPU: 8,
	MemoryMb: 8192,
	GPU: 0,
},
	"ecs.t5-c1m2.2xlarge":
{
	InstanceType: "ecs.t5-c1m2.2xlarge",
	VCPU: 8,
	MemoryMb: 16384,
	GPU: 0,
},
	"ecs.t5-c1m4.2xlarge":
{
	InstanceType: "ecs.t5-c1m4.2xlarge",
	VCPU: 8,
	MemoryMb: 32768,
	GPU: 0,
},
	"ecs.t5-c1m1.4xlarge":
{
	InstanceType: "ecs.t5-c1m1.4xlarge",
	VCPU: 16,
	MemoryMb: 16384,
	GPU: 0,
},
	"ecs.t5-c1m2.4xlarge":
{
	InstanceType: "ecs.t5-c1m2.4xlarge",
	VCPU: 16,
	MemoryMb: 32768,
	GPU: 0,
},
	"ecs.xn4.small":
{
	InstanceType: "ecs.xn4.small",
	VCPU: 1,
	MemoryMb: 1024,
	GPU: 0,
},
	"ecs.n4.small":
{
	InstanceType: "ecs.n4.small",
	VCPU: 1,
	MemoryMb: 2048,
	GPU: 0,
},
	"ecs.n4.large":
{
	InstanceType: "ecs.n4.large",
	VCPU: 2,
	MemoryMb: 4096,
	GPU: 0,
},
	"ecs.n4.xlarge":
{
	InstanceType: "ecs.n4.xlarge",
	VCPU: 4,
	MemoryMb: 8192,
	GPU: 0,
},
	"ecs.n4.2xlarge":
{
	InstanceType: "ecs.n4.2xlarge",
	VCPU: 8,
	MemoryMb: 16384,
	GPU: 0,
},
	"ecs.n4.4xlarge":
{
	InstanceType: "ecs.n4.4xlarge",
	VCPU: 16,
	MemoryMb: 32768,
	GPU: 0,
},
	"ecs.n4.8xlarge":
{
	InstanceType: "ecs.n4.8xlarge",
	VCPU: 32,
	MemoryMb: 65536,
	GPU: 0,
},
	"ecs.mn4.small":
{
	InstanceType: "ecs.mn4.small",
	VCPU: 1,
	MemoryMb: 4096,
	GPU: 0,
},
	"ecs.mn4.large":
{
	InstanceType: "ecs.mn4.large",
	VCPU: 2,
	MemoryMb: 8192,
	GPU: 0,
},
	"ecs.mn4.xlarge":
{
	InstanceType: "ecs.mn4.xlarge",
	VCPU: 4,
	MemoryMb: 16384,
	GPU: 0,
},
	"ecs.mn4.2xlarge":
{
	InstanceType: "ecs.mn4.2xlarge",
	VCPU: 8,
	MemoryMb: 32768,
	GPU: 0,
},
	"ecs.mn4.4xlarge":
{
	InstanceType: "ecs.mn4.4xlarge",
	VCPU: 16,
	MemoryMb: 65536,
	GPU: 0,
},
	"ecs.e4.small":
{
	InstanceType: "ecs.e4.small",
	VCPU: 1,
	MemoryMb: 8192,
	GPU: 0,
},
}
