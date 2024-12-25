package emailpermutator

import (
	"strings"
)

// Generate common email permutations based on first name, last name, and nickname.
func commonEmails(firstName, lastName, nickName string, domain string, emailOutput *[]string) {
	emailEnder := "@" + domain

	addEmail := func(email string) {
		*emailOutput = append(*emailOutput, email+emailEnder)
	}

	addEmail(firstName) // First Name
	if nickName != "" && nickName != firstName {
		addEmail(nickName) // Nick Name
	}
	addEmail(lastName)             // Last Name
	addEmail(firstName + lastName) // {fn}{ln}
	if nickName != "" && nickName != firstName {
		addEmail(nickName + lastName) // {nn}{ln}
	}
	addEmail(firstName + "." + lastName) // {fn}.{ln}
	if nickName != "" && nickName != firstName {
		addEmail(nickName + "." + lastName) // {nn}.{ln}
	}
	addEmail(string(firstName[0]) + lastName) // {fi}{ln}
	if nickName != "" && nickName != firstName {
		addEmail(string(nickName[0]) + lastName) // {ni}{ln}
	}
	addEmail(string(firstName[0]) + "." + lastName) // {fi}.{ln}
	if nickName != "" && nickName != firstName {
		addEmail(string(nickName[0]) + "." + lastName) // {ni}.{ln}
	}
	addEmail(firstName + string(lastName[0])) // {fn}{li}
	if nickName != "" && nickName != firstName {
		addEmail(nickName + string(lastName[0])) // {nn}{li}
	}
	addEmail(firstName + "." + string(lastName[0])) // {fn}.{li}
	if nickName != "" && nickName != firstName {
		addEmail(nickName + "." + string(lastName[0])) // {nn}.{li}
	}
	addEmail(string(firstName[0]) + string(lastName[0])) // {fi}{li}
	if nickName != "" && nickName != firstName {
		addEmail(string(nickName[0]) + string(lastName[0])) // {ni}{li}
	}
	addEmail(string(firstName[0]) + "." + string(lastName[0])) // {fi}.{li}
	if nickName != "" && nickName != firstName {
		addEmail(string(nickName[0]) + "." + string(lastName[0])) // {ni}.{li}
	}
}

// Generate less common email permutations based on first name, last name, and nickname.
func lessCommonEmails(firstName, lastName, nickName string, domain string, emailOutput *[]string) {
	emailEnder := "@" + domain

	addEmail := func(email string) {
		*emailOutput = append(*emailOutput, email+emailEnder)
	}

	addEmail(lastName + firstName) // {ln}{fn}
	if nickName != "" && nickName != firstName {
		addEmail(lastName + nickName) // {ln}{nn}
	}
	addEmail(lastName + "." + firstName) // {ln}.{fn}
	if nickName != "" && nickName != firstName {
		addEmail(lastName + "." + nickName) // {ln}.{nn}
	}
	addEmail(lastName + firstName[:1]) // {ln}{fi}
	if nickName != "" && nickName != firstName {
		addEmail(lastName + nickName[:1]) // {ln}{ni}
	}
	addEmail(lastName + "." + firstName[:1]) // {ln}.{fi}
	if nickName != "" && nickName != firstName {
		addEmail(lastName + "." + nickName[:1]) // {ln}.{ni}
	}
	addEmail(string(lastName[0]) + firstName) // {li}{fn}
	if nickName != "" && nickName != firstName {
		addEmail(string(lastName[0]) + nickName) // {li}{nn}
	}
	addEmail(string(lastName[0]) + "." + firstName) // {li}.{fn}
	if nickName != "" && nickName != firstName {
		addEmail(string(lastName[0]) + "." + nickName) // {li}.{nn}
	}
	addEmail(string(lastName[0]) + firstName[:1]) // {li}{fi}
	if nickName != "" && nickName != firstName {
		addEmail(string(lastName[0]) + nickName[:1]) // {li}{ni}
	}
	addEmail(string(lastName[0]) + "." + firstName[:1]) // {li}.{fi}
	if nickName != "" && nickName != firstName {
		addEmail(string(lastName[0]) + "." + nickName[:1]) // {li}.{ni}
	}
}

// Generate email permutations based on first name, last name, nickname, and middle name/initial.
func middleEmails(firstName, lastName, nickName, middleName string, domain string, emailOutput *[]string) {
	emailEnder := "@" + domain

	addEmail := func(email string) {
		*emailOutput = append(*emailOutput, email+emailEnder)
	}

	if middleName != "" {
		addEmail(firstName + middleName + lastName) // {fn}{mn}{ln}
		if nickName != "" && nickName != firstName {
			addEmail(nickName + middleName + lastName) // {nn}{mn}{ln}
		}
		addEmail(firstName + "." + middleName + "." + lastName) // {fn}.{mn}.{ln}
		if nickName != "" && nickName != firstName {
			addEmail(nickName + "." + middleName + "." + lastName) // {nn}.{mn}.{ln}
		}
	}
}

// Generate email permutations based on first name, last name, nickname, and middle name/initial with dashes.
func dashEmails(firstName, lastName, nickName, middleName string, domain string, emailOutput *[]string) {
	emailEnder := "@" + domain

	addEmail := func(email string) {
		*emailOutput = append(*emailOutput, email+emailEnder)
	}

	if middleName != "" {
		addEmail(firstName + "-" + middleName + "-" + lastName) // {fn}-{mn}-{ln}
		if nickName != "" && nickName != firstName {
			addEmail(nickName + "-" + middleName + "-" + lastName) // {nn}-{mn}-{ln}
		}
	}
}

// Generate email permutations based on first name, last name, nickname, and middle name/initial with underscores.
func underscoreEmails(firstName, lastName, nickName, middleName string, domain string, emailOutput *[]string) {
	emailEnder := "@" + domain

	addEmail := func(email string) {
		*emailOutput = append(*emailOutput, email+emailEnder)
	}

	if middleName != "" {
		addEmail(firstName + "_" + middleName + "_" + lastName) // {fn}_{mn}_{ln}
		if nickName != "" && nickName != firstName {
			addEmail(nickName + "" + middleName + "" + lastName) // {nn}{mn}{ln}
		}
	}
}

// Generate email permutations with the given data.
func Permute(data struct {
	FirstName, LastName, NickName, MiddleName, Domain1, Domain2, Domain3 string
}) []string {
	firstName := strings.ToLower(strings.TrimSpace(data.FirstName))
	lastName := strings.ToLower(strings.TrimSpace(data.LastName))
	nickName := strings.ToLower(strings.TrimSpace(data.NickName))
	middleName := strings.ToLower(strings.TrimSpace(data.MiddleName))
	domain1 := strings.ToLower(strings.TrimSpace(data.Domain1))
	domain2 := strings.ToLower(strings.TrimSpace(data.Domain2))
	domain3 := strings.ToLower(strings.TrimSpace(data.Domain3))

	var emailOutput []string

	commonEmails(firstName, lastName, nickName, domain1, &emailOutput)
	lessCommonEmails(firstName, lastName, nickName, domain1, &emailOutput)
	middleEmails(firstName, lastName, nickName, middleName, domain1, &emailOutput)
	dashEmails(firstName, lastName, nickName, middleName, domain1, &emailOutput)
	underscoreEmails(firstName, lastName, nickName, middleName, domain1, &emailOutput)

	if domain2 != "" {
		commonEmails(firstName, lastName, nickName, domain2, &emailOutput)
		lessCommonEmails(firstName, lastName, nickName, domain2, &emailOutput)
		middleEmails(firstName, lastName, nickName, middleName, domain2, &emailOutput)
		dashEmails(firstName, lastName, nickName, middleName, domain2, &emailOutput)
		underscoreEmails(firstName, lastName, nickName, middleName, domain2, &emailOutput)
	}

	if domain3 != "" {
		commonEmails(firstName, lastName, nickName, domain3, &emailOutput)
		lessCommonEmails(firstName, lastName, nickName, domain3, &emailOutput)
		middleEmails(firstName, lastName, nickName, middleName, domain3, &emailOutput)
		dashEmails(firstName, lastName, nickName, middleName, domain3, &emailOutput)
		underscoreEmails(firstName, lastName, nickName, middleName, domain3, &emailOutput)
	}

	return emailOutput
}
