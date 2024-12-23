# go-picturebook-sfomuseum

Go application to create a "picturebook" using the SFO Museum API.

## Description

## Creating a "picturebook" of items in your SFO Museum "shoebox"

## Tools

### picturebook

```
$> ./bin/picturebook -h
  -access-token string
    	A valid SFO Museum API access token to retrieve your shoebox items (must have "read" permissions).
  -bleed float
    	An additional bleed area to add (on all four sides) to the size of your picturebook.
  -border float
    	The size of the border around images. (default 0.01)
  -dpi float
    	The DPI (dots per inch) resolution for your picturebook. (default 150)
  -even-only
    	Only include images on even-numbered pages.
  -filename string
    	The filename (path) for your picturebook. (default "shoebox.pdf")
  -fill-page
    	If necessary rotate image 90 degrees to use the most available page space. Note that any '-process' flags involving colour space manipulation will automatically be applied to images after they have been rotated.
  -height float
    	A custom width to use as the size of your picturebook. Units are defined in inches by default. This flag overrides the -size flag when used in combination with the -width flag.
  -margin float
    	The margin around all sides of a page. If non-zero this value will be used to populate all the other -margin-(N) flags.
  -margin-bottom float
    	The margin around the bottom of each page. (default 1)
  -margin-left float
    	The margin around the left-hand side of each page. (default 1)
  -margin-right float
    	The margin around the right-hand side of each page. (default 1)
  -margin-top float
    	The margin around the top of each page. (default 1)
  -max-pages int
    	An optional value to indicate that a picturebook should not exceed this number of pages
  -odd-only
    	Only include images on odd-numbered pages.
  -orientation string
    	The orientation of your picturebook. Valid orientations are: 'P' and 'L' for portrait and landscape mode respectively. (default "P")
  -size string
    	A common paper size to use for the size of your picturebook. Valid sizes are: "a3", "a4", "a5", "letter", "legal", or "tabloid". (default "letter")
  -target-uri string
    	
  -units string
    	The unit of measurement to apply to the -height and -width flags. Valid options are inches, millimeters, centimeters (default "inches")
  -verbose
    	Display verbose output as the picturebook is created.
  -width float
    	A custom height to use as the size of your picturebook. Units are defined in inches by default. This flag overrides the -size flag when used in combination with the -height flag.
```

## Creating a SFO Museum API acccess token

The easiest and fastest way to create a SFO Museum API access token is to use the handy [Create a new access token for yourself](https://api.sfomuseum.org/oauth2/authenticate/like-magic/) webpage.

## See also

* https://github.com/aaronland/go-picturebook
* https://github.com/sfomuseum/go-sfomuseum-api
* https://you.sfomuseum.org
* https://api.sfomuseum.org