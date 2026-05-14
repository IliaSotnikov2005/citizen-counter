# citizen-counter

Консольная утилита на Go для подсчета частоты имен в текстовом файле.

## Возможности

- Подсчет количества повторений каждого имени в файле
- Настраиваемый размер буфера чтения файла
- Сортировка результатов по частоте или по алфавиту
- Вывод первых N результатов

## Установка

### Из исходников

```bash
git clone https://github.com/IliaSotnikov2005/citizen-counter.git
cd citizen-counter
go build -o citizen-counter ./cmd/
```
## Использование

###
```bash
# С параметрами по умолчанию
./citizen-counter file.txt

# Используя свои параметры
./citizen-counter -top=50 -buf=128 -sort=true file.txt
```

## Тестирование

```bash
go test -v ./...
```
