# Gateway_for_self_msg_broker


## Subscribing adding

To send the key to the system, you need to send an http request with JSON of the following type to the endpoint /add: {"key":"some_key","value":"some_value","topic":"some_topic"}. Authorization token in field "Authorization" is required in headers. If there are several topics on different nodes, the request will be sent to all nodes. If the topic is empty or does not exist, then look at the title ## Load balancing.

## YML

Do not forget to create a mb.yml file in the root of the application with the following contents:

    port: http_port
    prom_port: prometheus_port
    api_key: somekey
    wrong_topic: true/false
    nodes:
        - id: n
        topics:
            - topic1
            - topicn
        ip: some_ip:port
        scheme: http/s
        api_key: ""

If you specify any value in Scheme other than "https", it will be automatically set to "http". If an empty field is specified in node's api_key, the api_key from the first variable will be used. Every time the configuration is changed according to the ## Configuration header, all changes are immediately written to a file, which means that when the service is restarted, the configuration will always be up to date.

## Load balancing

If your application implies that the topic may be empty or it will not exist, then you need to set the wrong_topic variable to true. In this case, requests will be sent in a circle in order to all nodes existing in the configuration. If the topic is incorrect and the variable is false, the service will return a 204 Unauthorized response.

## Configuration

By endpoint /info you will get a json array with all the data in the nodes array from the yml configuration.
1. To add a node to the service, you need to send a json object in the request body: {"id":"n","addr":"ip:host","api_key":"","scheme":""} to endpoint /addnode.
2. To remove node from the service: {"id":"n"} in /rmnode.
3. To add a topic to the service and assign it to a specific node, you need to send a json object in the request body: {"topic":"some_topic","node_id":"existing_node_id"} to endpoint /addtopic.
4.  To remove topic from the service: {"topic":"some_topic"} in /rmtopic.
   
Authorization token is also required.

## Node setup

To use this gateway you need at least one working node of MessageBroker: https://github.com/MeDaLL1ST/MessageBroker . Make sure that the network is available before using the service.
