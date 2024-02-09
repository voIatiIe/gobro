import React from "react";
import { useState } from "react";
import { UrlInput } from "../UrlInput/UrlInput";
import styles from "./styles.module.css";


export const ToolBar = () => {
    const [URL, setURL] = useState<string>("");
    const [input, setInput] = useState<HTMLInputElement | null>(null);


    return (
        <div
            className={styles.toolBar}
        >
            <button>&#171;</button>
            <button>&#187;</button>
            <button>&#8635;</button>
            <UrlInput
                setUrl={setURL}
                setInput={setInput}
            />
            <button
                onClick={() => {setURL(""); if (input) {input.value = "";}}}
            >&#9587;</button>
            <button
                onClick={() => console.log(URL)}
            >
                Search
            </button>
        </div>
    )
}
