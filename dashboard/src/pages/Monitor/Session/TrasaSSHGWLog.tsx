import LinearProgress from '@material-ui/core/LinearProgress';
import React, { useRef } from 'react';
import { Terminal } from 'xterm';
// import * as fullscreen from 'xterm/lib/addons/fullscreen/fullscreen';
import '../../../utils/styles/xterm.css';
import { FitAddon } from 'xterm-addon-fit';

export default function TrasaSSHGWLog(props: any) {
  const termRef = useRef(new Terminal({ scrollback: Number.POSITIVE_INFINITY }));

  React.useEffect(() => {
    const term = termRef.current;
    const fitAddon = new FitAddon();

    term.loadAddon(fitAddon);
    fitAddon.fit();

    const container = document.getElementById('xterm');
    if (container) {
      term.open(container);
    }

    const d1 = props.sessionLog.replace(/\[3J/g, '');
    const d2 = d1.replace(/\[2J/g, '');
    const d3 = d2.replace(/\[H/g, '');
    const d4 = d3.replace(/\[3;J/g, '');

    term.write(d4);

    termRef.current = term;
    return () => {
      termRef.current.dispose();
    };
  }, [props.sessionLog]);

  return (
    <div>
      <div id="xterm">{props.TrasaSSHGWLog === '' ? <LinearProgress /> : null}</div>
    </div>
  );
}
