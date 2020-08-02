package main

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type Goods struct {
	Good_name        string
	Good_description string
	Good_price       float64
}

type UpdateGood struct {
	Is_description_update bool
	Is_price_update       bool
	Is_name_update        bool
}

func GetMd5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

/*
Функция для получения ключа ассоциативного массива.
Входную строку:
- переводим в нижний регистр
- убираем все, кроме цифр и букв из строки
- разбиваем строку на массив слов
- сортируем массив (для исключения зависимости порядка слов в наименовании товара)
- формируем строку из отсортированного массива и хешируем ее
*/
func getValidateKeyValue(str_orig string) string {

	if str_orig == "" {
		return ""
	}

	var result_str string
	var list_words_name []string

	result_str = strings.ToLower(str_orig)

	re_pattern := regexp.MustCompile(`[\d\wа-яА-Я]*`)
	re_str_res := re_pattern.FindAllString(result_str, -1)

	result_str = strings.Join(re_str_res, " ")

	list_words_name = strings.Fields(result_str)

	sort.Strings(list_words_name)

	result_str = strings.Trim(strings.Join(list_words_name, ""), "")

	result_str = GetMd5(result_str)

	return result_str
}
/*
Формируем список товаров
*/
func getListGoods(file_name string) map[string]Goods {
	file, err := os.Open(file_name)
	if err != nil {
		fmt.Println(err)
		fmt.Printf("file name: %q\n", file_name)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var good Goods
	var key_value string
	list_goods := make(map[string]Goods)

	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			fmt.Println("Error in scanner")
			continue
		}
		good.Good_name = ""
		good.Good_description = ""
		good.Good_price = 0

		if (len(strings.Split(scanner.Text(), ";"))) != 3 {
			fmt.Printf("Error! Structure of the string does not match expected: \"%s\"\n", scanner.Text())
			continue
		}

		good.Good_name = strings.Split(scanner.Text(), ";")[0]
		good.Good_description = strings.Split(scanner.Text(), ";")[1]
		good.Good_price, _ = strconv.ParseFloat(strings.Split(scanner.Text(), ";")[2], 64)
		key_value = getValidateKeyValue(good.Good_name)

		list_goods[key_value] = good

	}

	return list_goods
}
/*
Проверка на различие (тип string)
*/
func strCheckForDissimilarity(str1, str2 string) bool {
	return str1 != str2
}
/*
Проверка на различие (тип float)
*/
func floatCheckForDissimilarity(numb1, numb2 float64) bool {
	return numb1 != numb2
}

/*
Функция отображает изменения товарной позиции.
hex_key - ключ, по которому можно достать старое и новое значения из списков товаров
good - измененные атрибуты товарной позиции
old_list - список старых товаров
new_list - список новых товаров
name_attr - имя измененного атрибута
*/
func DisplayUpdates(hex_key string, good UpdateGood, old_list, new_list map[string]Goods, name_attr string) {
	fmt.Printf("\tUpdate %s:\n", name_attr)

	switch name_attr {
	case "name":
		if old_good, ok := old_list[hex_key]; ok {
			fmt.Printf("\t\told %s: %s\n", name_attr, old_good.Good_name)
		} else {
			fmt.Printf("Error! Not find %s in list_goods_old\n", hex_key)
		}

		if new_good, ok := new_list[hex_key]; ok {
			fmt.Printf("\t\tnew %s: %s\n", name_attr, new_good.Good_name)
		} else {
			fmt.Printf("Error! Not find %s in list_goods_new\n", hex_key)
		}
	case "description":
		if old_good, ok := old_list[hex_key]; ok {
			fmt.Printf("\t\told %s: %s\n", name_attr, old_good.Good_description)
		} else {
			fmt.Printf("Error! Not find %s in list_goods_old\n", hex_key)
		}

		if new_good, ok := new_list[hex_key]; ok {
			fmt.Printf("\t\tnew %s: %s\n", name_attr, new_good.Good_description)
		} else {
			fmt.Printf("Error! Not find %s in list_goods_new\n", hex_key)
		}
	case "price":
		if old_good, ok := old_list[hex_key]; ok {
			fmt.Printf("\t\told %s: %.2f\n", name_attr, old_good.Good_price)
		} else {
			fmt.Printf("Error! Not find %s in list_goods_old\n", hex_key)
		}

		if new_good, ok := new_list[hex_key]; ok {
			fmt.Printf("\t\tnew %s: %.2f\n", name_attr, new_good.Good_price)
		} else {
			fmt.Printf("Error! Not find %s in list_goods_new\n", hex_key)
		}
	}
}

func main() {

	path_to_file_old_goods := "old_list.csv"
	path_to_file_new_goods := "new_list.csv"

	list_goods_old := getListGoods(path_to_file_old_goods)
	list_goods_new := getListGoods(path_to_file_new_goods)

	//fmt.Println(list_goods_old)
	//fmt.Println(list_goods_new)

	fmt.Printf("Count old goods: %d\n", len(list_goods_old))
	fmt.Printf("Count new goods: %d\n", len(list_goods_new))

	list_update_goods := make(map[string]UpdateGood)
	list_treatment_new_goods := make(map[string]bool)

	var update_good UpdateGood
	var count_matches_goods int
	var list_deleted_goods []string
	for hex_key, old_good := range list_goods_old {

		update_good.Is_description_update = false
		update_good.Is_price_update = false
		update_good.Is_name_update = false

		//fmt.Println(old_key, old_good)
		if new_good, ok := list_goods_new[hex_key]; ok {
			count_matches_goods += 1

			list_treatment_new_goods[hex_key] = true

			update_good.Is_name_update = strCheckForDissimilarity(old_good.Good_name, new_good.Good_name)
			update_good.Is_description_update = strCheckForDissimilarity(old_good.Good_description, new_good.Good_description)
			update_good.Is_price_update = floatCheckForDissimilarity(old_good.Good_price, new_good.Good_price)

			if update_good.Is_name_update {
				fmt.Println("Warning. Names are synonymous, but not the same:")
				fmt.Printf("\tname from old list: %s\n", old_good.Good_name)
				fmt.Printf("\tname from new list: %s\n", new_good.Good_name)
				fmt.Println("\t---")
			}

			if update_good.Is_description_update || update_good.Is_price_update || update_good.Is_name_update {
				list_update_goods[hex_key] = update_good
			}

		} else {
			list_deleted_goods = append(list_deleted_goods, old_good.Good_name)
		}
	}
	fmt.Printf("Count of all matches product name: %d\n", count_matches_goods)

	if len(list_update_goods) == 0 {
		fmt.Println("================")
		fmt.Println("================")
		fmt.Println("================")
		fmt.Println("No products updates")
	} else {
		fmt.Println("================")
		fmt.Println("================")
		fmt.Println("================")
		fmt.Printf("Count of all goods updates: %d\n", len(list_update_goods))
		var list_attr_updates []string

		for hex_key, good := range list_update_goods {

			fmt.Printf("For good \"%s\":\n", list_goods_old[hex_key].Good_name)

			if good.Is_name_update {
				DisplayUpdates(hex_key, good, list_goods_old, list_goods_new, "name")
				list_attr_updates = append(list_attr_updates, "name")
			}

			if good.Is_description_update {
				DisplayUpdates(hex_key, good, list_goods_old, list_goods_new, "description")
				list_attr_updates = append(list_attr_updates, "description")
			}

			if good.Is_price_update {
				DisplayUpdates(hex_key, good, list_goods_old, list_goods_new, "price")
				list_attr_updates = append(list_attr_updates, "price")
			}

			fmt.Println("===")
		}
	}

	fmt.Println("================")
	fmt.Println("================")
	fmt.Println("================")
	fmt.Printf("Count of deleted goods: %d\n", len(list_deleted_goods))
	for _, name := range list_deleted_goods {
		fmt.Printf("\t%s\n", name)
	}

	var count_added_new_goods int
	var list_added_goods []string

	for hex_key, _ := range list_goods_new {
		if _, ok := list_treatment_new_goods[hex_key]; ok {
		} else {
			list_added_goods = append(list_added_goods, list_goods_new[hex_key].Good_name)
			count_added_new_goods += 1
		}
	}
	fmt.Println("================")
	fmt.Println("================")
	fmt.Println("================")
	fmt.Printf("Count of added new goods: %d\n", len(list_added_goods))
	for _, name := range list_added_goods {
		fmt.Printf("\t%s\n", name)
	}
}
