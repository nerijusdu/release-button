package audio

import (
	"fmt"
	"io"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/polly"
	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto"
)

type Synth struct {
	polly     *polly.Polly
	isEnabled bool
}

func NewSynth() *Synth {
	if os.Getenv("AWS_ACCESS_KEY_ID") == "" {
		return &Synth{isEnabled: false}
	}

	s := session.Must(session.NewSession())
	p := polly.New(s)

	return &Synth{
		polly:     p,
		isEnabled: true,
	}
}

func (s *Synth) Synthesize(text string) error {
	if !s.isEnabled {
		fmt.Println("AWS Polly is not enabled")
		return nil
	}

	output, err := s.polly.SynthesizeSpeech(&polly.SynthesizeSpeechInput{
		Text:         aws.String(text),
		OutputFormat: aws.String("mp3"),
		VoiceId:      aws.String("Stephen"),
		Engine:       aws.String("neural"),
	})

	if err != nil {
		return err
	}
	defer output.AudioStream.Close()

	d, err := mp3.NewDecoder(output.AudioStream)
	if err != nil {
		return err
	}

	c, err := oto.NewContext(d.SampleRate(), 2, 2, 8192)
	if err != nil {
		return err
	}
	defer c.Close()

	p := c.NewPlayer()
	defer p.Close()

	if _, err := io.Copy(p, d); err != nil {
		return err
	}

	return nil
}
