# Поисковик bsl модулей для sonar-scanner sonarqube

[![Release](https://img.shields.io/github/v/release/brobots-corporation/bsl2sonar.svg)](https://github.com/brobots-corporation/bsl2sonar/releases/latest)
[![lint and test](https://github.com/brobots-corporation/bsl2sonar/actions/workflows/lint-test.yml/badge.svg?branch=main)](https://github.com/brobots-corporation/bsl2sonar/actions/workflows/lint-test.yml)
[![codecov](https://codecov.io/gh/brobots-corporation/bsl2sonar/branch/main/graph/badge.svg?token=IZ00OLNPNN)](https://codecov.io/gh/brobots-corporation/bsl2sonar)
[![Go Report Card](https://goreportcard.com/badge/github.com/brobots-corporation/bsl2sonar)](https://goreportcard.com/report/github.com/brobots-corporation/bsl2sonar)
[![](https://img.shields.io/badge/license-GPL3-yellow.svg)](https://github.com/brobots-corporation/bsl2sonar/blob/main/LICENSE)

Поиск bsl файлов проекта (конфигурации 1С) по вхождению в подсистемы.

## Возможности

* Работа в ОС семейства: Linux, Windows, Mac OS X;
* Вывод полного или относительного пути к файлам с расширением .bsl;
* Вывод списка путей в файл sonar-project.properties или в поток стандартного вывода;
* Вывод кириллических символов в символах UNICODE;
* Генерация файла sonar-project.properties из шаблона.

## Сборка утилиты
* Скачать исходные файлы проекта, установить компилятор golang и собрать его командой:
  ```sh
    go build
    ```
## Установка и обновление
* Скачать бинарный файл bsl2sonar и расположить 
  его в необходимой для работы и хранения директории.
* Для обновления улиты необходимо скачать новую версии и заменить файл старой версии.
 
> Анализ файлов выгрузки выполняется для платформы 1С версии не ниже 8.3.10.

## Использование модуля

`bsl2sonar [-h] [-f FILE] [-a] [-u] [-v] [-l] [-g] srcdir parsephrases` - структура вызова утилиты

Обязательные аргументы:
* `srcdir` - путь к корневой папке с выгруженной конфигурацией 1с;
* `parsephrases` - префиксы подсистем, в которых будет осуществляться поиск путей до файлов объектов метаданных. Разделителем префиксов является пробел, к примеру `рн_ пк_ зс_`
  
Опциональные параметры:
* `-h, --help` - вызов справки;
* `-f FILE, --file FILE` - полный путь к файлу sonar-project.properties, в который будет выполняться выгрузка путей объектов метаданных на место переменной `$inclusions_line`;
* `-a, --absolute` - в случае указания флага будут выгружаться полные пути к файлам. Без флага только относительные пути;
* `-u, --unicode` - в случае указания флага будут выгружаться все кириллические символы в символах unicode;
* `-l, --logging` - в случае указания флага будут выводиться подробная информация;
* `-v, --version` - вывод версии скрипта;
* `-g, --generate` - генерация файла sonar-project.properties из шаблона;

Пример файла `sonar-project.properties` для первоначального запуска:

```properties
# Фильтры на включение в анализ. В примере ниже - только bsl и os файлы.
sonar.inclusions=$inclusions_line
```

### Пример использования скрипта в Linux

```sh
bsl2sonar "/Users/dummy/git/rn_erp/src/conf" "рн_ пс_" -u -f "/Users/dummy/git/rn_erp/sonar-project.properties"
```

### Пример использования скрипта в Windows

```cmd
bsl2sonar d:\rn_erp\src\conf  рн_ -u -f d:\rn_erp\sonar-project.properties
```