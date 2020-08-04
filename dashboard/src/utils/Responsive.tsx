// Responsive.js has utility functions to return size of elements based on window.innerHeight.
// There might be proper library to acheive this, but I am finding this easier way!! @sshahcodes.

// EchartElementHeight returns height for echartElement
export function EchartElementHeight() {
  // console.log('height: ', window.innerHeight, ' width: ', window.innerWidth)
  const h = window.innerHeight;
  switch (true) {
    case h < 750:
      return '160';
    case h > 750 && h < 820:
      return '190';
    case h > 820 && h < 970:
      return '200';
    default:
      // return '520'
      return '';
  }
}

export function EchartDivHeight() {
  const h = window.innerHeight;
  switch (true) {
    case h < 750:
      return { height: 220 };
    case h > 750 && h < 820:
      return { height: 250 };
    case h > 820 && h < 970:
      return { height: 270 };
    default:
      return { minHeight: 340 };
  }
}

export function EchartMapElementHeight() {
  const h = window.innerHeight;
  switch (true) {
    case h < 750:
      return '360';
    case h > 750 && h < 820:
      return '400';
    case h > 820 && h < 970:
      return '500';
    default:
      return '';
  }
}

export function EchartMapDivHeight() {
  const h = window.innerHeight;
  switch (true) {
    case h < 750:
      return { height: 520 };
    case h > 750 && h < 820:
      return { height: 570 };
    case h > 820 && h < 970:
      return { height: 620 };
    default:
      return { minHeight: 760 };
  }
}

export function HeaderFontSize() {
  const h = window.innerHeight;
  switch (true) {
    case h < 750:
      return '14px';
    case h > 750 && h <= 820:
      return '18px';
    case h > 820 && h < 970:
      return '18px';
    default:
      return '';
  }
}

export function TitleFontSize() {
  const h = window.innerHeight;
  switch (true) {
    case h < 750:
      return '21px';
    case h > 750 && h <= 820:
      return '22px';
    case h > 820 && h < 970:
      return '26px';
    default:
      return '';
  }
}

export function AGHeaderFontSize() {
  const h = window.innerHeight;
  switch (true) {
    case h < 750:
      return '14px';
    case h > 750 && h <= 820:
      return '16px';
    case h > 820 && h < 970:
      return '17px';
    default:
      return '';
  }
}

export function AGTitleFontSize() {
  const h = window.innerHeight;
  switch (true) {
    case h < 750:
      return '21px';
    case h > 750 && h <= 820:
      return '22px';
    case h > 820 && h < 970:
      return '28px';
    default:
      return '';
  }
}

export function Spacing() {
  const h = window.innerHeight;

  switch (true) {
    case h < 750:
      return 2;
    case h > 750 && h < 820:
      return 2;
    case h > 820 && h < 970:
      return 2;
    default:
      return 4;
  }
}

export function TitleTextFontSize() {
  const h = window.innerHeight;
  switch (true) {
    case h < 750:
      return '19px';
    case h > 750 && h <= 820:
      return '21px';
    case h > 820 && h < 970:
      return '26px';
    default:
      return '';
  }
}

export function ValueTextFontSize() {
  const h = window.innerHeight;
  switch (true) {
    case h < 750:
      return '19px';
    case h > 750 && h <= 820:
      return '21px';
    case h > 820 && h < 970:
      return '26px';
    default:
      return '';
  }
}
