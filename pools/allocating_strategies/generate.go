package pools

//go:generate echo ""
//go:generate echo "------> Loading allocation strategies as text into builtin_strategies.go for use at runtime "

//go:generate sh -c "echo \"package pools\" > ./builtin_strategies.go"
//go:generate sh -c "echo \"\" >> ./builtin_strategies.go"
//go:generate sh -c "echo \"const STRATEGY_INVOKER=\\``cat strategy_invoker.js`\\`\" >> ./builtin_strategies.go"
