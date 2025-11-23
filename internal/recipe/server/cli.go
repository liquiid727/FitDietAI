package server

import (
    "fmt"
    "os"
    "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use:   "recipe-agent",
    Short: "Recipe document assistant Agent",
}

func Execute() error {
    return rootCmd.Execute()
}

func init() {
    rootCmd.AddCommand(serveCmd)
    rootCmd.AddCommand(indexCmd)
}

var serveCmd = &cobra.Command{
    Use:   "serve",
    Short: "Start REST server",
    RunE: func(cmd *cobra.Command, args []string) error {
        return startHTTP()
    },
}

var indexCmd = &cobra.Command{
    Use:   "index",
    Short: "Index recipes from recipes/",
    RunE: func(cmd *cobra.Command, args []string) error {
        fmt.Fprintln(os.Stdout, "indexing recipes: TODO")
        return nil
    },
}