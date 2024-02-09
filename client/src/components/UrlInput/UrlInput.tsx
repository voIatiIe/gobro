import React from "react";
import styles from "./styles.module.css";

interface urlInputProps {
    setUrl: (url: string) => void
    setInput: (input: HTMLInputElement | null) => void
}

export const UrlInput: React.FC<urlInputProps> = ({ setUrl, setInput }) => {
    return (
        <div>
            <input
                className={styles.urlInput}
                type="text"
                onChange={(e) => setUrl(e.target.value)}
                ref={el => setInput(el)}
            />
        </div>
    )
}
