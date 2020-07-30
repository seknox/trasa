
export type TrasaHeaderAndContentProps = {
    children?: React.ReactNode;
   // history: History;
    tabHeaders: string[];
    pageName: pageProps[];
    Components: React.Component[];
    key: number;
}

type pageProps = {
    route: string;
    name: string;
}