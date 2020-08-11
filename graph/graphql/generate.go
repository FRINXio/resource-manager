package graphql

//go:generate echo ""
//go:generate echo "------> Generating graphql code from graph/graphql/schema"
//go:generate go run ./gqlgen.go

// Replace hardcoded symphony package import
//go:generate find . -name "tx_generated.go" -exec sed -i s/\"github.com\/facebookincubator\/symphony\/graph\/graphql\/generated\"/\"github.com\/net-auto\/resourceManager\/graph\/graphql\/generated\"/g {} +

// Remove generated Mutation() method from resolver, we need to override it with the one in resolver.go
//go:generate sed -i "s/func (r \\*Resolver) Mutation() .*/\\/\\/  Mutation() function removed in favour of resolver.go.Mutation()/g" resolver/schema.resolvers.go
