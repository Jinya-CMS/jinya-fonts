/* tslint:disable */
/* eslint-disable */
import { Designer } from '../models/designer';
import { Metadata } from '../models/metadata';
export interface Webfont {
  category: string;
  description?: string;
  designers?: Array<Designer>;
  fonts?: Array<Metadata>;
  googleFont?: boolean;
  license: string;
  name: string;
}
