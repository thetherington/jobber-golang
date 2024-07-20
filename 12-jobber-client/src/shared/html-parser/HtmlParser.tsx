import parse, { type DOMNode, Element, type HTMLReactParserOptions } from "html-react-parser";
import type { FC, ReactElement } from "react";

import type { IHtmlParserProps } from "../shared.interface";

const HtmlParser: FC<IHtmlParserProps> = ({ input }): ReactElement => {
  const parser = (input: string) => {
    const options: HTMLReactParserOptions = {
      replace: (domNode: DOMNode) => {
        if (domNode instanceof Element && domNode.firstChild && (domNode.firstChild as never)["name"] === "br") {
          return <></>;
        }
      },
    };
    return <>{parse(input, options)}</>;
  };

  return <>{parser(input)}</>;
};

export default HtmlParser;
