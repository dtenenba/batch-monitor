/*
This is meant to be a daemon that runs on the AMI that we use for Batch.
It monitors Docker events and sends them to an external data store.
Specifically it is looking for container start and destroy events.
On start we inspect the container and see if we can get the name of
the user who started the associated batch job by looking at
the container's environment variables. Saving the start and destroy events
should tell us how long the container was running. However, we also
need a lambda function that is triggered when instances are terminated
(or stopped?) so that if an instance is terminated abnormally before a
container can exit, we will still know when the instance went down.

What we are after is a breakdown of the EC2 instance hours by lab.
Challenges:

- What if multiple users from different labs used the same instance?
  How do we divide up the cost between these users?
- Sometimes instances stay on before or after jobs have run.
  So we need to track that time as well. What lab do we attribute
  that time to, esp if there are multiple users as in the previous point.
- Does ECS have any APIs that can help?
*/

package main

import (
	"context"
	"fmt"
	"os"

	"docker.io/go-docker"
	"docker.io/go-docker/api/types"
)

func main() {
	cli, err := docker.NewEnvClient()
	if err != nil {
		panic(err)
	}

	events, errchan := cli.Events(context.Background(), types.EventsOptions{})

	for {
		select {
		case event := <-events:
			fmt.Println("got an event")
			fmt.Println("status is", event.Status)
			if event.Status == "start" {
				fmt.Println("here is where we would inspect and get env vars")
			}
			fmt.Println(event)
		case <-errchan:
			fmt.Println("got a message on errchan, exiting")
			os.Exit(0)
		}
	}

	// containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	// if err != nil {
	// 	panic(err)
	// }
	//
	// for _, container := range containers {
	// 	fmt.Printf("%s %s\n", container.ID[:10], container.Image)
	// }
}
