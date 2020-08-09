package pkg

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v3"

	"os"
	"strings"
)

var (
	BastionsFilePath string
	homeN            string
)

func init() {
	homeN, _ := homedir.Dir()
	BastionsFilePath = homeN + "/.config/rodai/bastions.yml"

}

type Bastion struct {
	Name     string `yaml:"Name"`
	Username string `yaml:"Username"`
	Host     string `yaml:"Host"`
	Port     int    `yaml:"Port"`
	Key      string `yaml:"Key"`
}

type Bastions struct {
	Bastions []Bastion `yaml:"Bastions"`
}

func AddBastion(bastion Bastion) (string, int) {
	bastions := GetBastions(true)

	for _, f := range bastions.Bastions {
		if f.Name == bastion.Name {
			return "Failed creating bastion " + bastion.Name + " (Already exists).", 1
		}
	}

	bastions.Bastions = append(bastions.Bastions, bastion)

	CommitBastion(bastions)
	return "Success creating bastion " + bastion.Name, 0
}

func UpdateBastion(toUpdate Bastion, bastionName string) (string, int) {
	bastions := GetBastions(false)
	for _, bastion := range bastions.Bastions {
		if bastion.Name == bastionName {
			DeleteBastion(bastionName)
			bastion.Name = toUpdate.Name
			bastion.Host = toUpdate.Host
			bastion.Port = toUpdate.Port
			bastion.Key = toUpdate.Key
			AddBastion(bastion)
			return fmt.Sprintf("Update bastion %s (now %s)", bastionName, toUpdate.Name), 0
		}
	}

	return "Failed updating bastion " + bastionName + " (Not found)", 1
}

func DeleteBastion(toDelete string) (string, int) {

	bastions := GetBastions(false)
	for i, bastion := range bastions.Bastions {
		if bastion.Name == toDelete {
			bastions.Bastions = append(bastions.Bastions[:i], bastions.Bastions[i+1:]...)
			CommitBastion(bastions)
			return "Deleting bastion " + bastion.Name, 0
		}
	}
	return "Failed deleting bastion " + toDelete + " (Not found)", 1
}

func GetBastions(isAdd bool) Bastions {

	var bastions Bastions
	err := os.MkdirAll(ConfigPath, os.ModePerm)
	CheckErr(err)

	_, err = os.OpenFile(BastionsFilePath, os.O_RDONLY|os.O_CREATE, os.ModePerm)
	CheckErr(err)

	bastionContent := FileContentAsByteArray(BastionsFilePath)
	CheckErr(err)

	if len(bastionContent) == 0 && !isAdd {
		fmt.Println(ColoredStatus(16), "No bastions yet registered.")
		os.Exit(0)
	}

	err = yaml.Unmarshal(bastionContent, &bastions)
	CheckErr(err)

	return bastions
}

func GetBastionConfigDetails(name string) (Bastion, error) {
	var bastions Bastions

	_, err := os.OpenFile(BastionsFilePath, os.O_RDONLY|os.O_CREATE, 0644)
	CheckErr(err)

	bastionsContent := FileContentAsByteArray(BastionsFilePath)
	CheckErr(err)

	err = yaml.Unmarshal(bastionsContent, &bastions)
	CheckErr(err)

	for _, bastion := range bastions.Bastions {
		if bastion.Name == name {
			b, err := CheckSSHKey(bastion)
			CheckErr(err)
			return b, nil
		}
	}

	return Bastion{}, errors.New("bastion not found")
}

func CommitBastion(bastions Bastions) {
	var b bytes.Buffer
	encoder := yaml.NewEncoder(&b)
	encoder.SetIndent(2)

	err := encoder.Encode(&bastions)
	CheckErr(err)

	file, err := os.OpenFile(BastionsFilePath, os.O_WRONLY|os.O_TRUNC, 0644)
	CheckErr(err)
	defer file.Close()

	_, err = file.WriteString(b.String())
	CheckErr(err)
}

func CheckSSHKey(bastion Bastion) (Bastion, error) {

	if !strings.Contains(bastion.Key, ".ssh") {
		keyFilePath := home + "/.ssh/" + bastion.Key
		_, err := os.Stat(keyFilePath)
		if os.IsNotExist(err) {
			return Bastion{}, errors.New(fmt.Sprintf("SSH Key %s not found in PATH %s\n", bastion.Key, keyFilePath))
		} else {
			bastion.Key = keyFilePath
			return bastion, nil
		}
	}
	return Bastion{}, nil
}
