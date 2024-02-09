import React from 'react';
import { useState, useEffect, useRef } from 'react';
import './App.css';
import { Screen } from '../Screen/Screen';
import { GobroClient } from '../../client/client';


export const App = () => {
  const [src, setSrc] = useState<string>("");
  const [client, setClient] = useState<GobroClient>();
  const [player, setPlayer] = useState<HTMLImageElement | null>(null);

  useEffect(() => {
    if (!player) return;

    const getViewportParameters = () => {
      const rect = player.getBoundingClientRect();
  
      return { 
        width: rect.width,
        height: rect.height,
        top: rect.top,
        left: rect.left
      };
    }

    const client = new GobroClient(getViewportParameters, setSrc);
    setClient(client);
  }, [player])

  return <Screen src={src} client={client} setPlayer={setPlayer}/>;
}
