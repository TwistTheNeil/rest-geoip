interface MaxmindCountry {
  IsInEuropeanUnion: boolean;
  ISOCode?: string;
};

interface MaxmindCity {
  Names?: any; 
};

interface MaxmindLocation {
  AccuracyRadius: number;
  Latitude: number;
  Longitude: number;
  MetroCode: number;
  TimeZone: string;
};

interface MaxmindPostal {
  Code: string;
};

interface MaxmindTraits {
  IsAnonymousProxy: boolean;
  IsSatelliteProvider: boolean;
};

export interface MaxmindBackendResponse {
  Country: MaxmindCountry;
  City: MaxmindCity;
  Location: MaxmindLocation;
  Postal: MaxmindPostal;
  Traits: MaxmindTraits;
  IP: string;
};

export interface ConfigBackendResponse {
  MaptilerToken: string;
  AdminNotice: string;
};
