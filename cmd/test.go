This file is part of Foobar.

Foobar is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

Foobar is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with Foobar. If not, see <https://www.gnu.org/licenses/>.

package main

import (
	"context"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/galgotech/fermions-workflow-runtime/pkg/bus"
	"github.com/galgotech/fermions-workflow-runtime/pkg/log"
)

func main() {
	log := log.New("test")
	connector, err := bus.NewRedis(log, "redis://localhost:6379/0")
	if err != nil {
		panic(err)
	}

	event := cloudevents.NewEvent()
	event.SetID("id")
	event.SetSource("event-source")
	event.SetType("event-type")
	err = event.SetData("application/json", map[string]string{"key1": "value1", "key2": "value2"})
	if err != nil {
		panic(err)
	}

	data, err := event.MarshalJSON()
	if err != nil {
		panic(err)
	}

	// for {
	err = connector.Publish(context.Background(), event.Source(), data)
	if err != nil {
		panic(err)
	}
	// }
}
