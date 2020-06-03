package main

import(
 "fmt"
 "time"
 "math"
 //"image"
 "image/color"
 "github.com/fogleman/gg"
 "github.com/nfnt/resize"
 //"os"
 )
//---------------- constants---------------------------------------------------------------
const PIX_PER_MM=2.83463
var fWIDPAPER float64
var fHEIPAPER float64
var fMARGIN_DER,fMARGIN_IZQ float64
var fMARGIN_TOP, fMARGIN_BOTTOM float64
var bGregoriano,bDrawFest bool
var iIndexBeginYear int 
var iSizFontMes,iFontMes,iSizFontYear int
var fSepDay,fVentaja float64
var sPathImg,sPathFont,sPath1x1 string
//-----------------------------------------------------------------------------------------
func SetTypeCalendar(bTypeCal bool){
    // establece si es calendario gregoriano o juliano 
    bGregoriano=bTypeCal 
}
//-----------------------------------------------------------------------------------------
func SetDrawFestivo(bDrwFst bool){
    // establece si es calendario gregoriano o juliano 
    bDrawFest=bDrwFst 
}
//-----------------------------------------------------------------------------------------
func SetSizePaper(fWidpaper, fHeiPaper float64){
    // establece el tamaño del papel en milimetros
    fWIDPAPER=PIX_PER_MM*fWidpaper
    fHEIPAPER=PIX_PER_MM*fHeiPaper
}
//-----------------------------------------------------------------------------------------
func SetMargins(fTop, fBottom,fIzq,fDer float64){
    // establece las margenes del papel en milimetros
    fMARGIN_DER=PIX_PER_MM*fDer
    fMARGIN_IZQ=PIX_PER_MM*fIzq
    fMARGIN_TOP=PIX_PER_MM*fTop
    fMARGIN_BOTTOM=PIX_PER_MM*fBottom
}
//-----------------------------------------------------------------------------------------
func IsYearBisiesto(iYear int)bool{
    // decide si un año es bisiesto 
    var bBis bool
    if iYear%4==0{
        if iYear%100==0{
            if iYear%400==0{
                bBis=true
            }else{
                bBis=false 
            }
        }else{
            bBis=true
        }
    }else{
        bBis=false 
    }
    return bBis
}
//-----------------------------------------------------------------------------------------
func GetZellerIndex(iZDia,iZMes,iZYear int)int{
	// calcula el indice de zeller 
	var a,y,m,d,dia float64 
	var idia,iM,iDia int 
	dia=float64(iZDia)
	a=float64(14-iZMes)/12 
	y=float64(iZYear-int(a)) 
    iM=12*int(a)-2
	m=float64(iZMes)+float64(iM)
	//fmt.Printf("a= %.0f y= %.1f m= %.1f\n",a,y,m)
	if bGregoriano==true{
        d=dia+y+math.Floor(y/4)-math.Floor(y/100)+math.Floor(y/400)+float64((31*m)/12)
    }else{
        d=5+dia+y+math.Floor(y/4)+float64((31*m)/12)
    }
    idia=int(d)%7
    //fmt.Printf("d= %.0f idia= %d\n",d,idia)
    switch(idia){
        case 0:iDia=6
        case 1:iDia=0
        case 2:iDia=1
        case 3:iDia=2
        case 4:iDia=3
        case 5:iDia=4
        case 6:iDia=5
    }
    return iDia
}
//-----------------------------------------------------------------------------------------
func check(e error) {
    if e != nil {
        panic(e)
    }
}
//----------------------------------------------------------------------------------
func GetBinary(sText string)string{
	// convierte a binario un texto 
	var s,sBin string 
	for i:=0;i<len(sText);i++{
		iBin:=sText[i]
		s=fmt.Sprintf("%b",iBin)
		sBin+=s
	}
	return sBin
}
//----------------------------------------------------------------------------------
func DrawCodeBarClass(dc *gg.Context,sText string,x,y,fWid,fHei float64,iOpt int){
	// dibuja un codigo de barras clásico 
	var clr1,clr2 color.Color
	str:=GetBinary(sText)
	//fmt.Printf("%s\n",str)
	switch(iOpt){
		case 0:
			clr1=color.RGBA{255,255,255,255}
			clr2=color.RGBA{0,0,0,255}
		case 1:
			clr1=color.RGBA{255,0,0,255}
			clr2=color.RGBA{0,255,0,255}
		case 2:
			clr1=color.RGBA{255,0,0,255}
			clr2=color.RGBA{0,0,255,255}
		case 3:
			clr1=color.RGBA{0,255,0,255}
			clr2=color.RGBA{0,0,255,255}
		case 4:
			clr1=color.RGBA{255,0,0,255}
			clr2=color.RGBA{255,255,0,255}
		case 5:
			clr1=color.RGBA{0,255,0,255}
			clr2=color.RGBA{255,255,0,255}
		case 6:
			clr1=color.RGBA{168,168,168,255}
			clr2=color.RGBA{0,0,0,255}
	}
	for i:=0;i<len(str);i++{
		iBin:=str[i]
		dc.DrawRectangle(x,y,fWid,fHei)
		if iBin==49{
			dc.SetColor(clr2)
		}else{
			dc.SetColor(clr1)
		}
		dc.Fill()
		x+=fWid
	}
}
//-----------------------------------------------------------------------------------------
func TieneTilde(str string)(bool,int){
    var bTilde bool=false
    var ind int=0
    for i:=0;i<len(str);i++{
        if str[i]==195{
            bTilde=true
            ind++
        }
    }
    return bTilde,ind
}
//-----------------------------------------------------------------------------------------
func DrawTextRot(dc *gg.Context,x,y float64,sText string,fAngle float64,iFont,iSize int,clrText color.Color){
    // dibuja texto con cierto angulo 
    var s,s2 string
    var iLong int=0
    var vaX,vaY,hyp float64=0,0,0
    sAuxf:=fmt.Sprintf("%d.ttf",iFont)
	sFont:=sPathFont+sAuxf
	err:= dc.LoadFontFace(sFont,float64(iSize))
	if err != nil {
		 panic(err)
    }
    dc.SetColor(clrText)
    vaX=x
    vaY=y
    fAng:=-(fAngle*math.Pi)/180
    b,iNumTildes:=TieneTilde(sText)
    if b==true{
        iLong=len(sText)-iNumTildes
    }else{
        iLong=len(sText)
    }
    for i:=0;i<iLong;i++{
        s=string([]rune(sText)[i])
        if i==0{
            dc.DrawString(s,vaX,vaY)
            s2=s
        }else{
            w, h := dc.MeasureString(s2)
            if fAngle==0||fAngle==90||fAngle==180||fAngle==270||fAngle==360{
                hyp+=(w+1)
            }else{
                hyp+=(h+1)
            }
            vaX=x+hyp*math.Cos(fAng)
            vaY=y+hyp*math.Sin(fAng)
            dc.DrawString(s,vaX,vaY)
            s2=s
        }
    }
}
//-----------------------------------------------------------------------------------------
func DrawTextVert(dc *gg.Context,x,y float64,sText string,iFont,iSize int,clrText color.Color){
    // dibuja un texto vertical 
    var vaY float64=0
    var s string
    var iLong int=0
    sAuxf:=fmt.Sprintf("%d.ttf",iFont)
	sFont:=sPathFont+sAuxf
	err:= dc.LoadFontFace(sFont,float64(iSize))
	if err != nil {
		 panic(err)
    }
    dc.SetColor(clrText)
    vaY=y
    b,iNumTildes:=TieneTilde(sText)
    if b==true{
        iLong=len(sText)-iNumTildes
    }else{
        iLong=len(sText)
    }
    for i:=0;i<iLong;i++{
        s=string([]rune(sText)[i])
        dc.DrawString(s,x,vaY)
        vaY-=float64(iSize)
    }
}
//-----------------------------------------------------------------------------------------
func DrawText(dc *gg.Context,x,y float64,sText string,iFont,iSize int,clrText color.Color){
    // dibuja un texto horizontal corriente
    sAuxf:=fmt.Sprintf("%d.ttf",iFont)
	sFont:=sPathFont+sAuxf
	err:= dc.LoadFontFace(sFont,float64(iSize))
	if err != nil {
		 panic(err)
    }
    dc.SetColor(clrText)
    dc.DrawString(sText,x,y)
}
//-----------------------------------------------------------------------------------------
func DrawTextAlign(dc *gg.Context,y float64,sText string,iFont,iSize int,sOption string,clrText color.Color)float64{
    // dibuja un texto horizontal corriente
    var x float64=0
    sAuxf:=fmt.Sprintf("%d.ttf",iFont)
	sFont:=sPathFont+sAuxf
	err:= dc.LoadFontFace(sFont,float64(iSize))
	if err != nil {
		 panic(err)
    }
    w, h := dc.MeasureString(sText)
    dc.SetColor(clrText)
    if sOption=="left"{
        x=fMARGIN_IZQ
        dc.DrawString(sText,x,y)
    }else if sOption=="right"{
        x=fWIDPAPER-fMARGIN_DER-w
        dc.DrawString(sText,x,y)
    }else if sOption=="center"{
        x=(fWIDPAPER-w)/2.0
        dc.DrawString(sText,x,y)
    }else{
        x=fMARGIN_IZQ
        dc.DrawString(sText,x,y)
    }
    return h
}
//-----------------------------------------------------------------------------------------
func SetMesParamUser(iTamFontMes,iTypeFont,iSepCmt int){
    // establece los parametros de usuario
    // para dibujar el mes 
    iSizFontMes=iTamFontMes
    iSizFontYear=3*iSizFontMes
    iFontMes=iTypeFont
    fVentaja=float64(iSepCmt)*PIX_PER_MM
}
//-----------------------------------------------------------------------------------------
func DrawMargins(dc *gg.Context){
    // dibuja una margen al calendario 
    dc.SetLineWidth(2.0)
    dc.SetRGB255(0,0,0)
    dc.MoveTo(fMARGIN_IZQ,fMARGIN_TOP)
    dc.LineTo(fWIDPAPER-fMARGIN_DER,fMARGIN_TOP)
    dc.LineTo(fWIDPAPER-fMARGIN_DER,fHEIPAPER-fMARGIN_BOTTOM)
    dc.LineTo(fMARGIN_IZQ,fHEIPAPER-fMARGIN_BOTTOM)
    dc.LineTo(fMARGIN_IZQ,fMARGIN_TOP)
    dc.Stroke()
}
//-----------------------------------------------------------------------------------------
func GetNumDiasMes(iMes,iYear int)int{
    var iDays int=0
    var bBisiesto bool 
    bBisiesto=IsYearBisiesto(iYear)
    switch(iMes){
        case 0: // enero
            iDays=31
        case 1: // febrero
            if bBisiesto==true{
                iDays=29
            }else{
                iDays=28
            }
        case 2: // marzo
            iDays=31
        case 3: // abril
            iDays=30
        case 4: // mayo
            iDays=31
        case 5: // junio
            iDays=30
        case 6: // julio
            iDays=31
        case 7: // agosto
            iDays=31
        case 8: // septiembre
            iDays=30
        case 9: // octubre
            iDays=31
        case 10: // noviembre
            iDays=30
        case 11: // diciembre
            iDays=31
    }
    return iDays
}
//-----------------------------------------------------------------------------------------
func GetNameMes(iMes int)string{
    // devuelve el nombre del mes 
    var sMes string
    switch(iMes){
        case 0: // enero
            sMes="Enero"
        case 1: // febrero
            sMes="Febrero"
        case 2: // marzo
            sMes="Marzo"
        case 3: // abril
            sMes="Abril"
        case 4: // mayo
            sMes="Mayo"
        case 5: // junio
            sMes="Junio"
        case 6: // julio
            sMes="Julio"
        case 7: // agosto
            sMes="Agosto"
        case 8: // septiembre
            sMes="Septiembre"
        case 9: // octubre
            sMes="Octubre"
        case 10: // noviembre
            sMes="Noviembre"
        case 11: // diciembre
            sMes="Diciembre"
    }
    return sMes
}
//-----------------------------------------------------------------------------------------
func IsFestivoDay(iMes,iDay int)bool{
    // devuelve true si un dia es festivo 
    var bFest bool=false
    switch(iMes){
        case 0: // enero
            if iDay==1||iDay==6{
                bFest=true
            }
        case 1: // febrero
                bFest=false
        case 2: // marzo
            if iDay==23{
                bFest=true
            }
        case 3: // abril
            if iDay==9||iDay==10{
                bFest=true
            }
        case 4: // mayo
            if iDay==1||iDay==25{
                bFest=true
            }
        case 5: // junio
            if iDay==15||iDay==22||iDay==29{
                bFest=true
            }
        case 6: // julio
            if iDay==20{
                bFest=true
            }
        case 7: // agosto
            if iDay==7||iDay==17{
                bFest=true
            }
        case 8: // septiembre
                bFest=false
        case 9: // octubre
            if iDay==12{
                bFest=true
            }
        case 10: // noviembre
            if iDay==2||iDay==16{
                bFest=true
            }
        case 11: // diciembre
            if iDay==8||iDay==25{
                bFest=true
            }
    }
    return bFest
}
//-----------------------------------------------------------------------------------------
func GetFestivoText(iMes,iDay int)string{
    // devuelve true si un dia es festivo 
    var sFest string
    switch(iMes){
        case 0: // enero
            if iDay==1{
                sFest="Año Nuevo"
            }else if iDay==6{
                sFest="Día de los Reyes Magos"
            }
        
        case 1: // febrero
                sFest=""
        case 2: // marzo
            if iDay==23{
                sFest="Día de San José"
            }
        case 3: // abril
            if iDay==9{
                sFest="Jueves Santo"
            }else if iDay==10{
                sFest="Viernes Santo"
            }
        case 4: // mayo
            if iDay==1{
                sFest="Día del Trabajo"
            }else if iDay==25{
                sFest="Día de la Ascensión"
            }
        case 5: // junio
            if iDay==15{
                sFest="Corpus Christi"
            }else if iDay==22{
                sFest="Sagrado Corazón"
            }else if iDay==29{
                sFest="San Pedro y San Pablo"
            }
        case 6: // julio
            if iDay==20{
               sFest="Día de la Independencia"
            }
        case 7: // agosto
            if iDay==7{
                sFest="Batalla de Boyacá"
            }else if iDay==17{
                sFest="La asunción de la Virgen"
            }
        case 8: // septiembre
               sFest=""
        case 9: // octubre
            if iDay==12{
                sFest="Día de la Raza"
            }
        case 10: // noviembre
            if iDay==2{
                sFest="Todos los Santos"
            }else if iDay==16{
                sFest="Independencia de Cartagena"
            }
        case 11: // diciembre
            if iDay==8{
                sFest="Día de la Inmaculada Concepción"
            }else if iDay==25{
                sFest="Día de Navidad"
            }
    }
    return sFest
}
//-----------------------------------------------------------------------------------------
func GetNameDay(iDay int,iOption int)string{
    // devuelve la inicial del dia segun indice
    var sName string
    if iOption==0{
        switch(iDay){
        case 0:// lunes
            sName="Lunes"
        case 1:// lunes
            sName="Martes"
        case 2:// lunes
            sName="Miércoles"
        case 3:// lunes
            sName="Jueves"
        case 4:// lunes
            sName="Viernes"
        case 5:// lunes
            sName="Sábado"
        case 6:// lunes
            sName="Domingo"
        }
    }else{
        switch(iDay){
        case 0:// lunes
            sName="L"
        case 1:// lunes
            sName="M"
        case 2:// lunes
            sName="M"
        case 3:// lunes
            sName="J"
        case 4:// lunes
            sName="V"
        case 5:// lunes
            sName="S"
        case 6:// lunes
            sName="D"
        }
        
    }
    return sName
}
//-----------------------------------------------------------------------------------------
func DrawMes(dc *gg.Context,fX,fY float64,iMes,iYear int)(float64,float64){
    // dibuja un mes 
    var i,iBeginDay,iNumDays,iDay int=0,0,0,0
    var vaX,vaY,fWidMes,fHeiMes float64=0,0,0,0
    var sFont,sDay,sMes string
    iBeginDay=GetZellerIndex(1,iMes+1,iYear)
    iNumDays=GetNumDiasMes(iMes,iYear)
    sMes=GetNameMes(iMes)
    fmt.Printf("Mes: %s Begin Day: %d End day: %d\n",sMes,iBeginDay,iNumDays)
    // carga fuente para los titulos de los dias
    sAuxf:=fmt.Sprintf("%d.ttf",iFontMes)
	sFont=sPathFont+sAuxf
	if err := dc.LoadFontFace(sFont,float64(iSizFontMes)); err != nil {
		 panic(err)
    }
    fSepDay=float64(iSizFontMes+iSizFontMes/2.0)
    fWidMes=7*fSepDay
    // coloca el titulo del mes 
    vaX=fX+(fSepDay*8.0-float64(len(sMes)*iSizFontMes))/2.0
    vaY=fY
    clrTitle:=color.RGBA{0,0,128,255}
    dc.SetColor(clrTitle)
    dc.DrawString(sMes,vaX,vaY)
    vaX=fX
    vaY+=float64(iSizFontMes)
    // coloca los titulos para el día
    clrDayTit:=color.RGBA{0,0,255,255}
    dc.SetColor(clrDayTit)
    for i=0;i<7;i++{
        sDay=GetNameDay(i,1)
        dc.DrawString(sDay,vaX,vaY)  
        vaX+=fSepDay
    }
    vaY+=float64(iSizFontMes)
    clrDay:=color.RGBA{0,0,0,255}
    clrFiesta:=color.RGBA{255,0,0,255}
    vaX=fX+float64(iBeginDay)*fSepDay
    iDay=iBeginDay
    // dibuja los dias del mes
    for i=0;i<iNumDays;i++{
        bFest:=IsFestivoDay(iMes,i+1)
        if iDay%7==0{
            vaX=fX
            vaY+=float64(iSizFontMes)
            fHeiMes+=float64(iSizFontMes)
            iDay=0
        }
        if iDay%6==0&&iDay!=0{
            dc.SetColor(clrFiesta)
        }else if bFest==true{
            if bDrawFest==true{
                dc.SetColor(clrFiesta)
            }else{
                dc.SetColor(clrDay)
            }
        }else{
            dc.SetColor(clrDay)
        }
        sDay=fmt.Sprintf("%d",i+1)
        dc.DrawString(sDay,vaX,vaY)
        vaX+=fSepDay
        iDay++
    }
    return fWidMes,fHeiMes
}
//-----------------------------------------------------------------------------------------
func DrawMesInd(dc *gg.Context,fX,fY float64,iMes,iYear,iSizFont int)(float64,float64){
    // dibuja un mes 
    var i,iBeginDay,iNumDays,iDay int=0,0,0,0
    var vaX,vaY,fWidMes,fHeiMes float64=0,0,0,0
    var sFont,sDay,sMes string
    iBeginDay=GetZellerIndex(1,iMes+1,iYear)
    iNumDays=GetNumDiasMes(iMes,iYear)
    sMes=GetNameMes(iMes)
    fmt.Printf("Mes: %s Begin Day: %d End day: %d\n",sMes,iBeginDay,iNumDays)
    // carga fuente para los titulos de los dias
    sAuxf:=fmt.Sprintf("%d.ttf",iFontMes)
	sFont=sPathFont+sAuxf
	if err := dc.LoadFontFace(sFont,float64(iSizFont)); err != nil {
		 panic(err)
    }
    fSepDay=float64(iSizFont+iSizFont/2.0)
    fWidMes=7*fSepDay
    // coloca el titulo del mes 
    vaX=fX+(fSepDay*8.0-float64(len(sMes)*iSizFont))/2.0
    vaY=fY
    clrTitle:=color.RGBA{0,0,128,255}
    dc.SetColor(clrTitle)
    dc.DrawString(sMes,vaX,vaY)
    vaX=fX
    vaY+=float64(iSizFont)
    // coloca los titulos para el día
    clrDayTit:=color.RGBA{0,0,255,255}
    dc.SetColor(clrDayTit)
    for i=0;i<7;i++{
        sDay=GetNameDay(i,1)
        dc.DrawString(sDay,vaX,vaY)  
        vaX+=fSepDay
    }
    vaY+=float64(iSizFont)
    clrDay:=color.RGBA{0,0,0,255}
    clrFiesta:=color.RGBA{255,0,0,255}
    vaX=fX+float64(iBeginDay)*fSepDay
    iDay=iBeginDay
    // dibuja los dias del mes
    for i=0;i<iNumDays;i++{
        bFest:=IsFestivoDay(iMes,i+1)
        if iDay%7==0{
            vaX=fX
            vaY+=float64(iSizFont)
            fHeiMes+=float64(iSizFont)
            iDay=0
        }
        if iDay%6==0&&iDay!=0{
            dc.SetColor(clrFiesta)
        }else if bFest==true{
            if bDrawFest==true{
                dc.SetColor(clrFiesta)
            }else{
                dc.SetColor(clrDay)
            }
        }else{
            dc.SetColor(clrDay)
        }
        sDay=fmt.Sprintf("%d",i+1)
        dc.DrawString(sDay,vaX,vaY)
        vaX+=fSepDay
        iDay++
    }
    return fWidMes,fHeiMes
}
//-----------------------------------------------------------------------------------------
func DrawMesProg(dc *gg.Context,iIndMes,iYear int){
    // dibuja un mes de programador en (x,y) con wid y hei en mm
    var x,y,saltoH,saltoV float64
    var i,iSizFont int 
    saltoH=(fWIDPAPER-fMARGIN_IZQ-fMARGIN_DER)/7.0
    saltoV=(fHEIPAPER-fMARGIN_TOP-fMARGIN_BOTTOM-y)/7.0
    iSizFont=int(saltoV-saltoV/3.0)
    y=fMARGIN_TOP+float64(iSizFont)
    sMes:=GetNameMes(iIndMes)
	sTitMes:=fmt.Sprintf("%s / %d",sMes,iYear)
    _=DrawTextAlign(dc,y,sTitMes,64,iSizFont,"center",color.RGBA{0,0,0,255}) 
    x=fMARGIN_IZQ+float64(iSizFont/2.0)
    y+=saltoV
    // dibuja los días 
    for i=0;i<7;i++{
        sDay:=GetNameDay(i,1)
        if i<6{
            DrawText(dc,x,y,sDay,65,iSizFont,color.RGBA{0,0,0,255})  
        }else{
            DrawText(dc,x,y,sDay,65,iSizFont,color.RGBA{255,0,0,255})
        }
        x+=saltoH 
    }
    // dibuja la malla del mes 
    y+=float64(iSizFont/3.0)
    antY:=y
    x=fMARGIN_IZQ 
    dc.SetColor(color.RGBA{128,128,128,255})
    dc.SetLineWidth(2.0)
    // dibuja lineas verticales 
    for i=0;i<=7;i++{
        dc.MoveTo(x,y)
        dc.LineTo(x,y+5*saltoV)
        x+=saltoH 
    }
    dc.Stroke()
    // dibuja lineas horizontales 
    x=fMARGIN_IZQ
    for i=0;i<6;i++{
        dc.MoveTo(x,y)
        dc.LineTo(x+7*saltoH,y)
        y+=saltoV 
    }
    dc.Stroke()
    // dibuja los días del mes 
    iDay:=GetZellerIndex(1,iIndMes+1,iYear)
    iNumDias:=GetNumDiasMes(iIndMes,iYear)
    iSizFont=int(saltoV/2.0)
    y=antY+float64(iSizFont)
    x=fMARGIN_DER+float64(iSizFont/2.0)+float64(iDay)*saltoH
    var clr color.Color 
    clr=color.RGBA{128,128,128,255}
    for i=1;i<=iNumDias;i++{
        str:=fmt.Sprintf("%d",i)
		bFest:=IsFestivoDay(iIndMes,i)
        if iDay%6==0&&iDay!=0{
            clr=color.RGBA{255,0,0,255}
        }else if bFest==true{
            if bDrawFest==true{
                clr=color.RGBA{255,0,0,255}
            }else{
                clr=color.RGBA{128,128,128,255}
            }
        }else{
            clr=color.RGBA{128,128,128,255}
        }
        DrawText(dc,x,y,str,65,iSizFont,clr)
        x+=saltoH 
        iDay++
        if iDay%7==0{
           x=fMARGIN_DER+float64(iSizFont/2.0)
           y+=saltoV
		   iDay=0
        }
    }
	// dibuja calendarios pequeños 
	x=fMARGIN_IZQ
	y=fMARGIN_TOP
	for i=0;i<12;i++{
		w,h:=DrawMesInd(dc,x,y,i,iYear,8)
		x+=(w+w/5) 
		if i==5{
			x=fMARGIN_IZQ
			y+=2*h
		}
	}
	// dibuja los textos 
	y=fMARGIN_TOP+60
	x=fWIDPAPER/2.0+400
	DrawTextPublic(dc,x,y,20)
	// dibuja barra de código
	x=fMARGIN_IZQ
	y=fHEIPAPER-fMARGIN_BOTTOM
	str:="Copyright 2019 Horacio Useche"
	DrawCodeBarClass(dc,str,x,y,1,10,6)
}
//-----------------------------------------------------------------------------------------
func DrawTextPublic(dc *gg.Context,x,y float64,iSizFont int){
	// dibuja texto publicitario 
	clr:=color.RGBA{0,0,255,255}
	str:="ESU: La moneda virtual para el 2020"
	DrawText(dc,x,y,str,27,3*iSizFont,clr)
	clr=color.RGBA{64,64,64,255}
	y+=float64(iSizFont+5)
	str="Use ORIGIN para distinguir copias de documentos originales"
	DrawText(dc,x,y,str,62,iSizFont,clr)
	y+=float64(iSizFont)
	str="Use SCAM para burlar la seguridad de Google, Yahoo, Outlook, etc"
	DrawText(dc,x,y,str,62,iSizFont,clr)
	y+=float64(iSizFont)
	str="El Software de UseSoft33 no se distribuye por canales regulares!"
	DrawText(dc,x,y,str,62,iSizFont,clr)
	y+=float64(iSizFont)
	str="Ha sido vetado por Google Play y considerado peligroso!!!"
	DrawText(dc,x,y,str,62,iSizFont,clr)
	y+=float64(iSizFont)
	str="Más información en usesoft33@gmail.com"
	DrawText(dc,x,y,str,62,iSizFont,clr)
}
//-----------------------------------------------------------------------------------------
func GetWidMes()float64{
    // devuelve el ancho de un mes 
    var fWid float64=0
    fWid=float64(fSepDay*7.0)
    return fWid
}
//-----------------------------------------------------------------------------------------
func GetHeiMes()float64{
    // devuelve el ancho de un mes 
    var fHei float64=0
    fHei=float64(iSizFontMes*9.0)
    return fHei
}
//-----------------------------------------------------------------------------------------
func DrawYearText(dc *gg.Context,iYear int){
    // dibuja un año solo texto
    var vaX,vaY float64
    var fWidYear float64=0
    var sFont string
    iFontYear:=69
    sAuxf:=fmt.Sprintf("%d.ttf",iFontYear)
	sFont=sPathFont+sAuxf
	if err := dc.LoadFontFace(sFont,float64(iSizFontYear)); err != nil {
		 panic(err)
    }
    sYear:=fmt.Sprintf("%d",iYear)
    fWidYear=2*GetWidMes()+fVentaja
    fMid:=fWidYear/2.0    // la mitad
    vaX=fMARGIN_IZQ+4*fMid
    vaY=fMARGIN_TOP+float64(iSizFontYear)
    clrYear:=color.RGBA{0,0,255,255}
    dc.SetColor(clrYear)
    //dc.DrawString(sYear,vaX,vaY)
    DrawTextAlign(dc,vaY,sYear,iFontYear,iSizFontYear,"center",clrYear)
    vaY+=float64(iSizFontYear)
    vaX=fMARGIN_IZQ
    for i:=1;i<=12;i++{
         _,_=DrawMes(dc,vaX,vaY,i-1,iYear)
         vaX+=GetWidMes()+fVentaja
         if i%3==0{
             vaX=fMARGIN_IZQ
             vaY+=GetHeiMes()
         }
    }
    // dibuja codigo de barras 
    vaX=fMARGIN_IZQ
    vaY=fHEIPAPER-fMARGIN_BOTTOM
    str:="Copyright 2019 Horacio Useche Losada"
	DrawCodeBarClass(dc,str,vaX,vaY,1,10,6)
    DrawTextAlign(dc,vaY+8,str,iFontYear,8,"right",clrYear)
}
//-----------------------------------------------------------------------------------------
func DrawYearImg(dc *gg.Context,iYear int){
    // dibuja un año solo texto
    var vaX,vaY float64
    var fWidYear,fWidMes float64=0,0
    var sFont string
    iFontYear:=69
    sAuxf:=fmt.Sprintf("%d.ttf",iFontYear)
	sFont=sPathFont+sAuxf
	if err := dc.LoadFontFace(sFont,float64(iSizFontYear)); err != nil {
		 panic(err)
    }
    sYear:=fmt.Sprintf("%d",iYear)
    fWidYear=2*GetWidMes()+fVentaja
    fMid:=fWidYear/2.0    // la mitad
    vaX=fMARGIN_IZQ+4*fMid
    vaY=fMARGIN_TOP+float64(iSizFontYear)
    clrYear:=color.RGBA{0,0,255,255}
    dc.SetColor(clrYear)
    dc.DrawString(sYear,vaX,vaY)
    vaY+=float64(iSizFontYear)
    vaX=fMARGIN_IZQ
    iMes:=0
    ind:=0
    indImg:=0
    for i:=0;i<24;i++{
         if i==12{
             vaX=fMARGIN_IZQ+GetWidMes()+fVentaja
             vaY=fMARGIN_TOP+2.0*float64(iSizFontYear)
             ind=1
         }
         if(ind%2==0){
            fWidMes,_=DrawMes(dc,vaX,vaY,iMes,iYear)
            iMes++
         }else{
             //im image.Image 
             sAux:=fmt.Sprintf("%d.png",indImg)
			 sName:=sPathImg+sAux
             im, _ := gg.LoadImage(sName)
             m := resize.Resize(uint(fWidMes), 0, im, resize.Lanczos3)
             dc.DrawImage(m,int(vaX),int(vaY))
             indImg++
         }
         vaY+=GetHeiMes()
         ind++
    }
    DrawTextFestivos(dc,fMARGIN_IZQ,vaY,iYear)
    DrawTextCentral(dc,vaY)
    DrawMarkWater(dc)
    DrawCredits(dc)
}
//-----------------------------------------------------------------------------------------
func GetNumDiasFestivos(iYear int)int{
    // retorna el numero de dias festivos
    var bRec bool=false
    var iFest int=0
    for i:=0;i<11;i++{
        iDiasMes:=GetNumDiasMes(i,iYear)
        for j:=0;j<iDiasMes;j++{
            bRec=IsFestivoDay(i,j)
            if bRec==true{
                iFest++
            }
        }
    }
    return iFest
}
//-----------------------------------------------------------------------------------------
func DrawTextFestivos(dc *gg.Context,fX,fY float64,iYear int){
    // dibuja el texto con los festivos 
    var vaX,vaY,x float64=0,0,0
    var sFont,sText,sName string
    var bRec bool=false
    var ind int=0
    //var sF=make([]string,30) 
    iFontText:=63
    iSizFont:=iSizFontMes/3+5
    sAuxf:=fmt.Sprintf("%d.ttf",iFontText)
	sFont=sPathFont+sAuxf
	if err := dc.LoadFontFace(sFont,float64(iSizFont)); err != nil {
		 panic(err)
    }
    vaX=fX
    vaY=fY
    sTitle:="Días Festivos"
    clrTitle:=color.RGBA{0,0,0,255}
    dc.SetColor(clrTitle)
    dc.DrawString(sTitle,vaX,vaY)
    fSalto:=fWIDPAPER/3.0
    vaY+=float64(iSizFont)
    clrText:=color.RGBA{0,0,0,255}
    // dibuja el texto de los días festivos 
    for i:=0;i<12;i++{
        iDiasMes:=GetNumDiasMes(i,iYear)
        for j:=0;j<iDiasMes;j++{
            bRec=IsFestivoDay(i,j)
            if bRec==true{
                if ind<6{
                    vaX=fX
                    x=vaX+float64(6*iSizFont)
                }else if ind==6{
                    vaX=fX+fSalto
                    vaY=fY+float64(iSizFont)
                    x=vaX+float64(6*iSizFont)
                }else if ind==12{
                    vaX=fX+2*fSalto
                    vaY=fY+float64(iSizFont)
                    x=vaX+float64(8*iSizFont)
                }
                sText=GetFestivoText(i,j)
                sName=GetNameMes(i)
                sAux:=fmt.Sprintf("%d de %s:",j,sName)
                fmt.Printf("%s %s\n",sAux,sText)
                dc.SetColor(color.RGBA{255,0,0,255})
                dc.DrawString(sAux,vaX,vaY)
                dc.SetColor(clrText)
                dc.DrawString(sText,x,vaY)
                vaY+=float64(iSizFont)
                ind++
            }
        }
    }
}
//-----------------------------------------------------------------------------------------
func DrawTextCentral(dc *gg.Context,fY float64){
    // dibuja el texto al centro del calendario 
    x:=GetWidMes()+fVentaja/2.0+60.0
    y:=fY-fVentaja/2-120.0
    str:="Feliz Navidad y Próspero Año Nuevo"
    iFont:=63
    iSize:=72
    clrText:=color.RGBA{128,128,128,255}
    DrawTextVert(dc,x,y,str,iFont,iSize,clrText)
}
//-----------------------------------------------------------------------------------------
func DrawMarkWater(dc *gg.Context){
    // dibuja texto como marcas de agua 
    iFont:=71
    iSize:=10
    clrText:=color.RGBA{230,230,230,245}
    str:=" UseSoft33 - usesoft33@gmail.com "
    x:=5.0
    y:=fHEIPAPER
    for i:=0;i<1;i++{
        DrawTextVert(dc,x,y,str,iFont,iSize,clrText)
        DrawTextVert(dc,x+float64(iSize+iSize/2),y,str,iFont,iSize,clrText)
        y-=float64(len(str)*iSize)
    }
    str="Copyright 2020 by Horacio Useche Losada, usesoft33@gmail.com"
    iSize=10
    y=fHEIPAPER-float64(iSize)
    _=DrawTextAlign(dc,y,str,63,iSize,"right",color.RGBA{0,0,0,245})
    str="Imagen Volcán Nevado del Huila 2019"
    y=fHEIPAPER-float64(3*iSize)
    _=DrawTextAlign(dc,y,str,63,iSize,"right",color.RGBA{64,64,64,245})
	 str="Visit our blog https://usesoft33.home.blog/2019/12/16/origin/"
    y=fHEIPAPER-float64(2*iSize)
    _=DrawTextAlign(dc,y,str,63,iSize,"right",color.RGBA{128,128,128,245})
}
//-----------------------------------------------------------------------------------------
func DrawCredits(dc *gg.Context){
    // dibuja los créditos del fabricante 
    iFont:=23
    iSize:=10
    clrText:=color.RGBA{128,128,128,255}
    x:=fMARGIN_IZQ
    y:=fHEIPAPER-float64(6*iSize)
    str:="ORIGIN: software to protect digital documents and distinguish copies of their original !!!"
    DrawText(dc,x,y,str,iFont,iSize,clrText)
	y+=float64(iSize)
    str="SCAM: software to fool computer systems and filter prohibited files, such as files binaries"
    DrawText(dc,x,y,str,iFont,iSize,clrText)
    y+=float64(iSize)
    str="UseMath: High precision library to make math software with C language, performance and speed in your routines"
    DrawText(dc,x,y,str,iFont,iSize,clrText)
    y+=float64(iSize)
    str="UseGis: Geographic information system to make thematic maps oriented to work with geographical coordinates"
    DrawText(dc,x,y,str,iFont,iSize,clrText)
    y+=float64(iSize)
    str="UseFractal: Software to make giant and amazing fractal images with an unpublished format"
    DrawText(dc,x,y,str,iFont,iSize,clrText)
    y+=float64(iSize)
    str="Ask for it through usesoft33@gmail.com. We will send you a disk using the services of certified mail (DHL)"
    DrawText(dc,x,y,str,iFont,iSize,clrText)
}
//-----------------------------------------------------------------------------------------
func DrawYearImg3x4(dc *gg.Context,iYear int){
    // dibuja un año con graficos
    // organizando 4 filas
    var vaX,vaY float64
    var fWidMes float64=0
    vaY=fMARGIN_TOP+10.0
    sYear:=fmt.Sprintf("%d",iYear)
    h:=DrawTextAlign(dc,vaY,sYear,64,60,"center",color.RGBA{0,0,0,255})
    vaY+=h
    str:="Merry Christmas and Happy New Year"
    //str:="Feliz Navidad y Próspero Año"
    DrawTextAlign(dc,vaY,str,62,50,"center",color.RGBA{128,128,128,255})
    vaY+=h
    str="Are the wishes of SCAM and ORIGIN App"
    //str="Les desea familia Pueto Useche"
    h=DrawTextAlign(dc,vaY,str,63,20,"center",color.RGBA{0,0,0,255})
    vaY+=h
    str="The software banned by google because everything good is dangerous"
    h=DrawTextAlign(dc,vaY,str,63,20,"center",color.RGBA{0,0,0,255})
    vaY+=h
    str="Visit our blog at https://wordpress.com/view/usesoft33.home.blog"
    h=DrawTextAlign(dc,vaY,str,63,16,"center",color.RGBA{0,0,0,255})
    vaY+=2*h
    vaX=fMARGIN_IZQ
    iMes:=0
    ind:=0
    indImg:=0
    for i:=0;i<24;i++{
         if i==6||i==12||i==18{
             vaX=fMARGIN_IZQ
             vaY+=GetHeiMes()
         }
         if(ind%2==0){
            fWidMes,_=DrawMes(dc,vaX,vaY,iMes,iYear)
            iMes++
         }else{
             //im image.Image
             sAux:=fmt.Sprintf("%d.png",indImg)
			 sName:=sPathImg+sAux
             im, _ := gg.LoadImage(sName)
             m := resize.Resize(uint(fWidMes), 0, im, resize.Lanczos3)
             dc.DrawImage(m,int(vaX),int(vaY))
             indImg++
         }
         vaX+=GetWidMes()
         ind++
    }
    vaY+=GetHeiMes()
	str="Copyright 2020 Horacio Useche"
	DrawCodeBarClass(dc,str,fMARGIN_IZQ,905,1,10,6)
    DrawTextFestivos(dc,fMARGIN_IZQ,vaY,iYear)
    DrawMarkWater(dc)
    DrawCredits(dc)
}
//------------------------------------------------------------------------------------------------
func DrawYearImg1x1(dc *gg.Context,iYear int){
    // dibuja un año con graficos
    // organizando 4 filas
    var vaX,vaY float64
    vaY=fMARGIN_TOP+20.0
    sYear:=fmt.Sprintf("%d",iYear)
    h:=DrawTextAlign(dc,vaY,sYear,64,90,"center",color.RGBA{0,0,0,255})
    vaY+=h
    str:="Merry Christmas and Happy New Year"
    DrawTextAlign(dc,vaY,str,62,60,"center",color.RGBA{128,128,128,255})
    vaY+=h
    str="Are the wishes of ORIGIN & SCAM"
    h=DrawTextAlign(dc,vaY,str,63,40,"center",color.RGBA{0,0,0,255})
    vaY+=h
    str="The software banned by google because everything good is dangerous"
    h=DrawTextAlign(dc,vaY,str,63,26,"center",color.RGBA{0,0,0,255})
    vaX=fMARGIN_IZQ
    vaY+=h
    // dibuja la imagen
    fWidImg:=fWIDPAPER-fMARGIN_IZQ-fMARGIN_DER
    vaX=fMARGIN_IZQ
    im, _ := gg.LoadImage(sPath1x1)
    m := resize.Resize(uint(fWidImg), 0, im, resize.Lanczos3)
    dc.DrawImage(m,int(vaX),int(vaY))
    // dibuja el calendario
    fLuz:=10.0
    vaY=fWidImg
    vaX=fMARGIN_IZQ+fLuz
    for i:=0;i<12;i++{
        if i==4||i==8{
             vaX=fMARGIN_IZQ
             vaY+=GetHeiMes()
         }
         _,_=DrawMes(dc,vaX,vaY,i,iYear)
         vaX+=(GetWidMes()+1.6*fLuz)
    }
    vaY+=GetHeiMes()
    str="Copyright 2019 Horacio Useche"
	DrawCodeBarClass(dc,str,fMARGIN_IZQ,1550,1,10,6)
    DrawTextFestivos(dc,fMARGIN_IZQ,vaY,iYear)
    DrawMarkWater(dc)
    DrawCredits(dc)
}
//------------------------------------------------------------------------------------------------
func BuildYear(iYear int){
    // construye un calendario anual 
    SetSizePaper(360,580)
    dc := gg.NewContext(int(fWIDPAPER), int(fHEIPAPER))
    dc.SetRGB(1, 1, 1)
    dc.Clear() 
    SetMargins(20,20,20,20)
    SetMesParamUser(20,62,8)
    DrawYearImg1x1(dc,iYear)
    //DrawYearImg3x4(dc,iYear)
    //DrawYearImg(dc,iYear)
    //DrawYearText(dc,iYear)
    //DrawMes(dc,100.0,200.0,3)
    //fmt.Printf("Mes empieza en: %d\n",iRec)
    dc.SavePNG("./imgs/calend_year.png")
}
//------------------------------------------------------------------------------------------------
func BuildProgMes(iIndMes,iYear int,fHei float64){
    // construye un programador mensual 
    fWid:=(16.0/9.0)*fHei
    SetSizePaper(fWid,fHei)
    dc := gg.NewContext(int(fWIDPAPER), int(fHEIPAPER))
    dc.SetRGB(1, 1, 1)
    dc.Clear() 
    SetMargins(5,5,5,5)
    SetMesParamUser(40,62,30)
    DrawMesProg(dc,iIndMes,iYear)
    dc.SavePNG("./imgs/mesprog.png")
}
//------------------------------------------------------------------------------------------------
func ShowBeginYear(iYear, iNum int){
	// muestra los días que empieza y termina un año 
	var iYears int=iYear 
	for i:=0;i<=iNum;i++{
		iRec1:=GetZellerIndex(1,1,iYears)
		iRec2:=GetZellerIndex(31,12,iYears)
		str1:=GetNameDay(iRec1,0)
		str2:=GetNameDay(iRec2,0)
        fmt.Printf("%d : %s : %s\n",iYears,str1,str2)
		iYears++
	}
}
//------------------------------------------------------------------------------------------------
func main(){   
 str:=fmt.Sprintf("Tiempo antes de ejecutar: %v",time.Now())   
 fmt.Printf("Ejecutando UseCalend 2020 ...\n")
 sPathImg="/home/xun33/dev/devgo/projects/calend/imgs/aves/img_"
 sPathFont="/home/xun33/dev/Fonts/Serial/font_"
 sPath1x1="/home/xun33/dev/devgo/projects/calend/imgs/singles/nevado_huila_2.jpg"
 SetTypeCalendar(true)
 SetDrawFestivo(true)
 BuildYear(2020)
 //BuildProgMes(0,2020,400)
 //iRec:=GetZellerIndex(2,8,1953)
 //fmt.Printf("Zeller: %d\n",iRec)
 //ShowBeginYear(2013,12)
 fmt.Println(str)
 fmt.Printf("Tiempo después de ejecutar: %v \n",time.Now())
}



