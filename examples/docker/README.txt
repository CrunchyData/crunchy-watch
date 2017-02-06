this example assumes you have run the crunchy-containers
example/docker/master-replica example.  That example will
create a master container and replica container that form
a postgres cluster. 

This watch example looks for a master container, if not found, 
it will trigger a failover onto the replica container.

To test this example, delete the master pod and examine
the log of the watch container to see if perform the failover logic.
