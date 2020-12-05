# BrigolageCMS
A clone of BricolageCMS in Go.

# The Backstory
[BricolageCMS](https://bricolagecms.org) is a CMS [written in Perl](https://github.com/bricoleurs/bricolage).
I was interested in using it, but was never able to get it up and running.

I'm sort of halfsies-converting it Go to help me with both my Perl and Go web server skills.
And to see if I can make it easier to install.

# Copyrights and Licenses
My code, which is derived from the
[original](https://github.com/bricoleurs/bricolage),
is licensed under the 3-clause BSD license.
Please see the [LICENSE](LICENSE) file.

The BricolageCMS code is copyright by About.com, Kineticode, Inc., and many others.
It is licensed under the [3-clause BSD license](https://opensource.org/licenses/BSD-3-Clause).

    Copyright (c) 2001-2002 About.com.

    Changes Copyright (c) 2002-2009 Kineticode, Inc. and others.

    All rights reserved.
    
    Redistribution and use in source and binary forms, with or without modification,
    are permitted provided that the following conditions are met:
    
    * Redistributions of source code must retain the above copyright notice, this
      list of conditions and the following disclaimer.
    * Redistributions in binary form must reproduce the above copyright notice,
      this list of conditions and the following disclaimer in the documentation
      and/or other materials provided with the distribution.
    * Neither the name of About.com nor the names of its contributors may be used
      to endorse or promote products derived from this software without specific
      prior written permission.
    
        THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"  
        AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE  
        IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE  
        DISCLAIMED. IN NO EVENT SHALL THE REGENTS OR CONTRIBUTORS BE LIABLE FOR ANY  
        DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES  
        (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;  
        LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON  
        ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT  
        (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS  
        SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.  

Please see
[TheBricolageLicense](cmd/about/static/TheBricolageLicense)
in this repository for a full list of the BricolageCMS contributors.

## Website
There is a copy of the BricolageCMS website in the `cmd\about\public\` directory.
That website is copyright by the Bricolage Development Team and is licensed under the
[Creative Commons Attribution-NonCommercial 2.0 Generic (CC BY-NC 2.0) license](https://creativecommons.org/licenses/by-nc/2.0/legalcode).

The copy was created the evening of 2020/12/04 by using `wget` (see the `cmd\about\spider.sh` script) to mirror the site to local files.
I'm sure it includes too much (and too little), but it's helpful for me during the conversion.

## Icons
BricolageCMS uses icon sets that are covered by separate copyrights.

* The Nuvola icon set,
from [David Vignoni](http://icon-king.com),
is licensed under the
[GNU Lesser General Public License, version 2.1](https://www.gnu.org/licenses/old-licenses/lgpl-2.1.en.html):

    This library is free software; you can redistribute it and/or
    modify it under the terms of the GNU Lesser General Public
    License as published by the Free Software Foundation; either
    version 2.1 of the License, or (at your option) any later version.

    This library is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
    Lesser General Public License for more details.

    You should have received a copy of the GNU Lesser General Public
    License along with this library; if not, write to the Free Software
    Foundation, Inc., 59 Temple Place, Suite 330, Boston, MA  02111-1307  USA

* The Tango icon set,
from the [Tango Desktop Project](http://tango.freedesktop.org),
is licensed under the
[Creative Commons Attribution Share-Alike license](http://creativecommons.org/licenses/by-sa/2.5/).


