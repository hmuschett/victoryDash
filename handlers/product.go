package handlers

import (
	"strconv"
	"strings"
	"victorydash/configs"

	goshopify "github.com/bold-commerce/go-shopify"
)

//InsertProduct from LineItem on order into product_order
func InsertProduct(orderID int64, product goshopify.LineItem) {
	query := "INSERT vic.product_order SET  order_id=?, sku=?, vendor=?, quantity=?"
	quantity := product.Quantity * getQuantiOfPack(product.VariantTitle)
	configs.VicExec(query, orderID, product.SKU, product.Vendor, quantity)
}

func getQuantiOfPack(str string) (result int) {
	out := strings.Split(str, " ")
	value, _ := strconv.Atoi(out[0])

	return value
}
