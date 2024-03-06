/* tslint:disable */
/* eslint-disable */
import { HttpClient, HttpContext } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';

import { BaseService } from '../base-service';
import { ApiConfiguration } from '../api-configuration';
import { StrictHttpResponse } from '../strict-http-response';

import { Designer } from '../models/designer';
import { getFontByName } from '../fn/operations/get-font-by-name';
import { GetFontByName$Params } from '../fn/operations/get-font-by-name';
import { getFontDesignersByName } from '../fn/operations/get-font-designers-by-name';
import { GetFontDesignersByName$Params } from '../fn/operations/get-font-designers-by-name';
import { getFontFilesByName } from '../fn/operations/get-font-files-by-name';
import { GetFontFilesByName$Params } from '../fn/operations/get-font-files-by-name';
import { getFonts } from '../fn/operations/get-fonts';
import { GetFonts$Params } from '../fn/operations/get-fonts';
import { Metadata } from '../models/metadata';
import { Webfont } from '../models/webfont';

@Injectable({ providedIn: 'root' })
export class ApiService extends BaseService {
  constructor(config: ApiConfiguration, http: HttpClient) {
    super(config, http);
  }

  /** Path part for operation `getFonts()` */
  static readonly GetFontsPath = '/api/font';

  /**
   * Gets all fonts.
   *
   *
   *
   * This method provides access to the full `HttpResponse`, allowing access to response headers.
   * To access only the response body, use `getFonts()` instead.
   *
   * This method doesn't expect any request body.
   */
  getFonts$Response(params?: GetFonts$Params, context?: HttpContext): Observable<StrictHttpResponse<Array<Webfont>>> {
    return getFonts(this.http, this.rootUrl, params, context);
  }

  /**
   * Gets all fonts.
   *
   *
   *
   * This method provides access only to the response body.
   * To access the full response (for headers, for example), `getFonts$Response()` instead.
   *
   * This method doesn't expect any request body.
   */
  getFonts(params?: GetFonts$Params, context?: HttpContext): Observable<Array<Webfont>> {
    return this.getFonts$Response(params, context).pipe(
      map((r: StrictHttpResponse<Array<Webfont>>): Array<Webfont> => r.body)
    );
  }

  /** Path part for operation `getFontByName()` */
  static readonly GetFontByNamePath = '/api/font/{fontName}';

  /**
   * Gets the given font by name.
   *
   *
   *
   * This method provides access to the full `HttpResponse`, allowing access to response headers.
   * To access only the response body, use `getFontByName()` instead.
   *
   * This method doesn't expect any request body.
   */
  getFontByName$Response(params: GetFontByName$Params, context?: HttpContext): Observable<StrictHttpResponse<Webfont>> {
    return getFontByName(this.http, this.rootUrl, params, context);
  }

  /**
   * Gets the given font by name.
   *
   *
   *
   * This method provides access only to the response body.
   * To access the full response (for headers, for example), `getFontByName$Response()` instead.
   *
   * This method doesn't expect any request body.
   */
  getFontByName(params: GetFontByName$Params, context?: HttpContext): Observable<Webfont> {
    return this.getFontByName$Response(params, context).pipe(map((r: StrictHttpResponse<Webfont>): Webfont => r.body));
  }

  /** Path part for operation `getFontFilesByName()` */
  static readonly GetFontFilesByNamePath = '/api/font/{fontName}/file';

  /**
   * Gets the files for the given font.
   *
   *
   *
   * This method provides access to the full `HttpResponse`, allowing access to response headers.
   * To access only the response body, use `getFontFilesByName()` instead.
   *
   * This method doesn't expect any request body.
   */
  getFontFilesByName$Response(
    params: GetFontFilesByName$Params,
    context?: HttpContext
  ): Observable<StrictHttpResponse<Metadata>> {
    return getFontFilesByName(this.http, this.rootUrl, params, context);
  }

  /**
   * Gets the files for the given font.
   *
   *
   *
   * This method provides access only to the response body.
   * To access the full response (for headers, for example), `getFontFilesByName$Response()` instead.
   *
   * This method doesn't expect any request body.
   */
  getFontFilesByName(params: GetFontFilesByName$Params, context?: HttpContext): Observable<Metadata> {
    return this.getFontFilesByName$Response(params, context).pipe(
      map((r: StrictHttpResponse<Metadata>): Metadata => r.body)
    );
  }

  /** Path part for operation `getFontDesignersByName()` */
  static readonly GetFontDesignersByNamePath = '/api/font/{fontName}/designer';

  /**
   * Gets the designers for the given font.
   *
   *
   *
   * This method provides access to the full `HttpResponse`, allowing access to response headers.
   * To access only the response body, use `getFontDesignersByName()` instead.
   *
   * This method doesn't expect any request body.
   */
  getFontDesignersByName$Response(
    params: GetFontDesignersByName$Params,
    context?: HttpContext
  ): Observable<StrictHttpResponse<Designer>> {
    return getFontDesignersByName(this.http, this.rootUrl, params, context);
  }

  /**
   * Gets the designers for the given font.
   *
   *
   *
   * This method provides access only to the response body.
   * To access the full response (for headers, for example), `getFontDesignersByName$Response()` instead.
   *
   * This method doesn't expect any request body.
   */
  getFontDesignersByName(params: GetFontDesignersByName$Params, context?: HttpContext): Observable<Designer> {
    return this.getFontDesignersByName$Response(params, context).pipe(
      map((r: StrictHttpResponse<Designer>): Designer => r.body)
    );
  }
}
