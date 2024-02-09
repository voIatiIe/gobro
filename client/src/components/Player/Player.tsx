import React from "react";
import { useState } from "react";
import { useEffect } from "react";
import { GobroClient } from "../../client/client";
import styles from "./styles.module.css";


interface playerProps {
    src: string,
    client?: GobroClient,
    setPlayer: (client: HTMLImageElement | null) => void,
}

export const Player: React.FC<playerProps> = ({ src, client, setPlayer }) => {
    return (
        <img
            className={styles.player}
            src={src}
            onMouseMove={(e) => client?.onMouseMove(e)}
            onClick={(e) => client?.onClick(e)}
            ref={(img) => {setPlayer(img)}}
        />
    );
}
