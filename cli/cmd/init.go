package cmd

import (
	"errors"
	"fmt"
	"github.com/azarc-io/verathread-dev-toolkit/cli/internal/types"
	"github.com/azarc-io/verathread-dev-toolkit/cli/internal/util"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"io/fs"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

/************************************************************************/
/* TYPES
/************************************************************************/

type (
	InitCmd struct {
		projectName string
		description string
		title       string
		pwd         string
		privatePort string
		publicPort  string
		webPort     string
		debugPort   string
	}

	initProgram struct {
		progress progress.Model
		spinner  spinner.Model
		queue    []*processInitFile
		percent  float64
		total    int
		cmd      *InitCmd
		width    int
		height   int
		done     bool
		index    int
	}

	processInitFile struct {
		path string
	}
)

/************************************************************************/
/* COMMAND
/************************************************************************/

func (i *InitCmd) Cmd(cmd *cobra.Command, args []string) error {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	i.pwd = pwd
	i.projectName = filepath.Base(pwd)
	i.publicPort = "6020"
	i.privatePort = "6021"
	i.webPort = "3001"
	i.debugPort = "40001"

	if err := huh.NewInput().
		Title("What is the project name?").
		Description("Must be lower case with optional hyphens, preferably the same name as the git repository").
		Suggestions([]string{filepath.Base(pwd)}).
		Value(&i.projectName).
		Validate(func(s string) error {
			if s == "" {
				return errors.New("application name is required")
			}
			return nil
		}).Run(); err != nil {
		return err
	}

	if err := huh.NewInput().
		Title("Provide a short title for this Application?").
		Value(&i.title).
		Validate(func(s string) error {
			if s == "" {
				return errors.New("application title is required")
			}
			return nil
		}).Run(); err != nil {
		return err
	}

	if err := huh.NewInput().
		Title("Provide a description for this Application?").
		Value(&i.description).
		Validate(func(s string) error {
			if s == "" {
				return errors.New("application description is required")
			}
			return nil
		}).Run(); err != nil {
		return err
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Private Port").
				Value(&i.privatePort).
				Description("local private network port").
				CharLimit(4).
				Validate(i.validatePort(4, "private")),
		),
		huh.NewGroup(
			huh.NewInput().
				Title("Public Port").
				Value(&i.publicPort).
				Description("local public network port").
				CharLimit(4).
				Validate(i.validatePort(4, "public")),
		),
		huh.NewGroup(
			huh.NewInput().
				Title("Web Port").
				Value(&i.webPort).
				Description("web server network port").
				CharLimit(4).
				Validate(i.validatePort(4, "web server")),
		),
		huh.NewGroup(
			huh.NewInput().
				Title("Debug Port").
				Value(&i.debugPort).
				Description("backend debug port").
				CharLimit(5).
				Validate(i.validatePort(5, "debug")),
		),
	).WithLayout(huh.LayoutGrid(1, 3))

	if err := form.Run(); err != nil {
		return err
	}

	return i.runProgram()
}

/************************************************************************/
/* COMMAND HELPERS
/************************************************************************/

func (i *InitCmd) validatePort(minLength int, name string) func(port string) error {
	return func(port string) error {
		if len(port) < 4 {
			return fmt.Errorf("invalid port format, must contain %d digits", minLength)
		} else if port == "" {
			return fmt.Errorf("%s port is required", name)
		} else if _, err := strconv.Atoi(port); err != nil {
			return errors.New("invalid port format, must contain only numbers")
		}
		return nil
	}
}

func (i *InitCmd) runProgram() error {
	var (
		queue []*processInitFile
	)

	if err := filepath.WalkDir(i.pwd, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() && !strings.Contains(path, "node_modules") {
			if strings.HasSuffix(path, ".go") ||
				strings.HasSuffix(path, ".ts") ||
				strings.HasSuffix(path, ".html") ||
				strings.HasSuffix(path, ".md") ||
				strings.HasSuffix(path, ".yaml") ||
				strings.HasSuffix(path, ".yml") ||
				strings.HasSuffix(path, ".json") ||
				strings.HasSuffix(path, "Tiltfile") ||
				strings.Contains(path, "Dockerfile") {
				queue = append(queue, &processInitFile{
					path: path,
				})
			}
		}
		return nil
	}); err != nil {
		return err
	}

	if _, err := tea.NewProgram(&initProgram{
		progress: progress.New(
			progress.WithScaledGradient("#FF7CCB", "#FDFF8C"),
			progress.WithoutPercentage(),
			progress.WithWidth(80),
		),
		spinner: spinner.New(spinner.WithSpinner(spinner.Dot)),
		queue:   queue,
		total:   len(queue),
		cmd:     i,
	}).Run(); err != nil {
		return err
	}

	return nil
}

/************************************************************************/
/* COMMAND FACTORY
/************************************************************************/

func NewInitCmd() *InitCmd {
	return &InitCmd{}
}

/************************************************************************/
/* PROGRAM
/************************************************************************/

// Init is the first function that will be called. It returns an optional
// initial command. To not perform an initial command return nil.
func (i *initProgram) Init() tea.Cmd {
	return tea.Batch(i.nextCmd(), i.spinner.Tick)
}

// nextCmd get the next command to process
func (i *initProgram) nextCmd() tea.Cmd {
	//return func() tea.Msg {
	//	v, _ := i.queue.Dequeue()
	//	return v
	//}

	d := time.Millisecond * time.Duration(rand.Intn(100))
	return tea.Tick(d, func(t time.Time) tea.Msg {
		return i.queue[i.index]
	})
}

// Update is called when a message is received. Use it to inspect messages
// and, in response, update the model and/or send a command.
func (i *initProgram) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case *processInitFile:
		if i.index >= len(i.queue)-1 {
			i.done = true
			return i, tea.Sequence(
				tea.Printf("%s %s", types.CheckMark, msg.path),
				tea.Quit,
			)
		}

		f, err := util.ReadFile(msg.path)
		if err != nil {
			log.Error("error while reading file", "path", msg.path)
			return nil, tea.Quit
		}

		out, err := util.ParseTemplate(f, map[string]interface{}{
			"PROJECT_NAME":        i.cmd.projectName,
			"PROJECT_DESCRIPTION": i.cmd.description,
			"PROJECT_TITLE":       i.cmd.title,
			"PRIVATE_PORT":        i.cmd.privatePort,
			"PUBLIC_PORT":         i.cmd.publicPort,
			"WEB_PORT":            i.cmd.webPort,
			"DEBUG_PORT":          i.cmd.debugPort,
		})
		if err != nil {
			log.Error("error while parsing file", "path", msg.path)
			return nil, tea.Quit
		}

		if err := util.WriteFile([]byte(out), msg.path); err != nil {
			log.Error("error while parsing file", "path", msg.path)
			return nil, tea.Quit
		}

		i.index++
		i.percent = (float64(i.index) / float64(i.total)) * 1
		progressCmd := i.progress.SetPercent(i.percent)

		return i, tea.Batch(
			progressCmd,
			tea.Printf("%s %s", types.CheckMark, msg.path),
			i.nextCmd(),
		)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return i, tea.Quit
		}
	case spinner.TickMsg:
		var cmd tea.Cmd
		i.spinner, cmd = i.spinner.Update(msg)
		return i, cmd
	case progress.FrameMsg:
		newModel, cmd := i.progress.Update(msg)
		if newModel, ok := newModel.(progress.Model); ok {
			i.progress = newModel
		}
		return i, cmd
	case tea.WindowSizeMsg:
		i.width, i.height = msg.Width, msg.Height
	}

	return i, nil
}

// View renders the program's UI, which is just a string. The view is
// rendered after every Update.
func (i *initProgram) View() string {
	n := i.total
	w := lipgloss.Width(fmt.Sprintf("%d", n))

	if i.done {
		return types.DoneStyle.Render(fmt.Sprintf("Done! Processed %d files.\n", n))
	}

	pkgCount := fmt.Sprintf(" %*d/%*d", w, i.index, w, n)
	spin := i.spinner.View() + " "
	prog := i.progress.View()
	cellsAvail := max(0, i.width-lipgloss.Width(spin+pkgCount))
	pkgName := types.CurrentPkgNameStyle.Render(i.queue[i.index].path)
	info := lipgloss.NewStyle().MaxWidth(cellsAvail).Render("Processing " + pkgName)
	cellsRemaining := max(0, i.width-lipgloss.Width(spin+info+prog+pkgCount))
	gap := strings.Repeat(" ", cellsRemaining)

	return spin + info + gap + pkgCount + "\n" + prog
}
