package main

import (
    "fmt"
    "github.com/Chef1991/IssueTracker/models"
    "bufio"
    "os"
    "strings"
)

func main() {
    var (
        name string
        creatorId int
        description string
        shortDescription string
    )
    fmt.Println("init")
    reader := bufio.NewReader(os.Stdin)
    models.Init()
    fmt.Println("Enter a Project name:")
    name, _ = reader.ReadString('\n')
    fmt.Println("Enter Creator's id:")
    fmt.Scanf("%d\n", &creatorId)
    fmt.Println("Enter a Description:")
    description, _ = reader.ReadString('\n')
    fmt.Println("Enter a short descriptoin")
    shortDescription, _ = reader.ReadString('\n')
    name = strings.TrimSpace(name)
    description = strings.TrimSpace(description)
    shortDescription = strings.TrimSpace(shortDescription)
    creator, err := models.GetUserFromId(creatorId)
    if err != nil {
        fmt.Println(err)
        return
    }
    project, err := models.CreateProject(name, creator, description, shortDescription)
    if err != nil {
        fmt.Println("ERROR:", err)
    }
    fmt.Println(project)
}
