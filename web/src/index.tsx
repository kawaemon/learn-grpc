import React from "react";
import ReactDOM from "react-dom";

import { grpc } from "@improbable-eng/grpc-web";
import { Chat, ChatChat } from "./grpc/protocol_pb_service";
import { ChatMessage } from "./grpc/protocol_pb";

function Main() {
    const [messages, setMessages] = React.useState<string[]>([]);

    const [client, setClient] = React.useState<grpc.Client<
        ChatMessage,
        ChatMessage
    > | null>(null);

    React.useEffect(() => {
        console.log("recreated");

        const client = grpc.client<ChatMessage, ChatMessage, ChatChat>(
            Chat.Chat,
            {
                host: "https://lc.kawaemon.dev:3000",
                transport: grpc.WebsocketTransport()
            }
        );

        client.onHeaders(() => {
            setMessages((m) => [...m, "== connected! =="]);
        });

        client.onMessage((msg: ChatMessage) => {
            setMessages((m) => [...m, `server: ${msg.getMessage()}`]);
        });

        client.onEnd((code) => {
            setMessages((m) => [...m, `== disconnected(${code})! ==`]);
        });

        client.start(new grpc.Metadata());

        setClient(client);

        // return () => client.finishSend();
    }, [])

    const [input, setInput] = React.useState<string>("");

    return (
        <>
            <div>
                <input
                    onChange={(e) => setInput(e.target.value)}
                    placeholder="Type your message here!"
                />
                <button
                    disabled={client == null}
                    onClick={() => {
                        const message = new ChatMessage();
                        message.setMessage(input);
                        setMessages([...messages, `you: ${input}`]);
                        client!.send(message);
                    }}
                >
                    Send
                </button>
            </div>

            <div>
                {messages.map((v, i) => (
                    <div key={i}>{v}</div>
                ))}
            </div>
        </>
    );
}

ReactDOM.render(<Main />, document.getElementById("root"));
