package model

type Category struct {
	BaseModel
	Name             string `gorm:"type:varchar(20);not null;comment:商品名称"`
	ParentCategoryID int32
	ParentCategory   *Category
	Level            int32 `gorm:"type:int;not null;default:1;comment:商品分类级别"`
	IsTab            bool  `gorm:"type:bool;default:false;not null;comment:能否显示在tab栏"`
}

/*
先创建 商品分类表（有三级分类） 和 品牌分类表

同一个目录下的包名一定要一致
	同一个包名下面的东西可以直接用，不用import

Name 的类型不能为null
	实际开发过程中尽量设置为 not null

https://zhuanlan.zhihu.com/p/73997266
总结：
NULL 本身是一个特殊值，MySQL 采用特殊的方法来处理 NULL 值。从理解肉眼判断，操作符运算等操作上，可能和我们预期的效果不一致。
可能会给我们项目上的操作不符合预期。你必须要使用 IS NULL / IS NOT NULL 这种与普通 SQL 大相径庭的方式去处理 NULL。
尽管在存储空间上，在索引性能上可能并不比空值差，但是为了避免其身上特殊性，给项目带来不确定因素，因此建议默认值不要使用 NULL。

Level 商品分类级别
	用int类型还是int32类型,用int32，因为字段不长，尽量减少储存空间
	统一定义为int32，防止频繁类型转换，因为proto文件中没有int类型，只有int32、int64、uint32、uint64

ParentCategoryID 是外键字段，数据库中真正存储的
ParentCategory 是外键
gorm中如果想要自己指向自己，要使用指针

ParentCategoryID: 父分类的 ID，类型为 int32，并创建索引以提高查询性能。
ParentCategory: 指向父分类的指针，通过 ParentCategoryID 作为外键关联。通过这个字段，GORM 可以自动建立与 ParentCategoryID 的关联关系。

*/
