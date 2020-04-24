package tone

import (
	"fmt"
	"math"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
)

func createTone(sampleRate beep.SampleRate, freq float64) beep.Streamer {
	var playbackPos int
	return beep.StreamerFunc(func(samples [][2]float64) (n int, ok bool) {
		for i := range samples {
			amp := math.Sin(2.0 * math.Pi * freq / float64(sampleRate.N(time.Second)) * float64(playbackPos))
			samples[i][0] = amp
			samples[i][1] = amp
			playbackPos++
		}
		return len(samples), true
	})
}

type Note struct {
	Freq   float64
	Length int
	Delay  int
	Output string
}

func Play(notes []Note) {
	sr := beep.SampleRate(44100)
	speaker.Init(sr, sr.N(time.Millisecond*200))

	var seq []beep.Streamer
	done := make(chan struct{})
	for _, note := range notes {
		fmt.Println(note.Output)
		seq = append(seq, beep.Take(sr.N(time.Duration(note.Length)*time.Millisecond), createTone(sr, note.Freq)))
		seq = append(seq, beep.Silence(sr.N(time.Duration(note.Delay)*time.Millisecond)))
	}
	seq = append(seq, beep.Callback(func() {
		done <- struct{}{}
	}))
	speaker.Play(beep.Seq(seq...))
	<-done
}
