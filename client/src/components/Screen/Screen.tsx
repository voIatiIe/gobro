import React from "react";
import { ToolBar } from "../ToolBar/ToolBar";
import { Player } from "../Player/Player";
import { GobroClient } from "../../client/client";


interface screenProps {
    src: string,
    client?: GobroClient,
    setPlayer: (client: HTMLImageElement | null) => void,
}

export const Screen: React.FC<screenProps> = ({ src, client, setPlayer }) => {
    return (
        <>
            <ToolBar></ToolBar>
            <Player src={src} client={client} setPlayer={setPlayer}></Player>
        </>
    );
}
